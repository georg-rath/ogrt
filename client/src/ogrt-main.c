#include "ogrt-main.h"

/** global variables */
static bool  __ogrt_active   =  0;
static int   __daemon_socket = -1;
static pid_t __pid           =  0;
static pid_t __parent_pid    =  0;
static char  __hostname[HOST_NAME_MAX+1];
static uuid_t __uuid = {0};
static void (*saved_signal_handlers[32]) (int) = { NULL };

FILE *ogrt_log_file;
int ogrt_log_level;

/** real functions without hooks */
static int (*real_sigaction)(int, const struct sigaction *, struct sigaction *) = NULL;
static sighandler_t (*real_signal)(int, sighandler_t) = NULL;

/**
 * Initialize preload library.
 * Checks if tracing is activated using the environment variable OGRT_ACTIVE.
 * If it is: connect to the unix socket of the daemon. If establishing a connection to
 * the daemon fails the init function will return with a non-zero exit code, but program
 * execution will continue as normal.
 */
ON_INIT
int ogrt_preload_init_hook()
{
  ogrt_log_file = stderr;
  ogrt_log_level = OGRT_LOG_INFO;

  if(ogrt_env_enabled("OGRT_DBG")) {
    ogrt_log_level = OGRT_LOG_DBG;
  }
  Log(OGRT_LOG_DBG, "called init hook\n");

  if(ogrt_env_enabled("OGRT_SILENT") || ogrt_env_enabled("OGRT_SCHLEICHFAHRT")) {
    ogrt_log_level = OGRT_LOG_NOTHING;
  }

  if(ogrt_env_enabled("OGRT_DEBUG_INFO")) {
    cmdline_parser_print_version();
    printf("  OGRT_NET_HOST=%s\n  OGRT_NET_PORT=%s\n  OGRT_ENV_JOBID=%s\n  OGRT_ELF_SECTION_NAME=%s\n  OGRT_ELF_NOTE_TYPE=0x%x\n",
              OGRT_NET_HOST,      OGRT_NET_PORT,      OGRT_ENV_JOBID,      OGRT_ELF_SECTION_NAME,      OGRT_ELF_NOTE_TYPE);
    printf("  OGRT_MSG_SEND_USERNAME=%d\n  OGRT_MSG_SEND_HOSTNAME=%d\n  OGRT_MSG_SEND_ENVIRONMENT=%d\n",
              OGRT_MSG_SEND_USERNAME,      OGRT_MSG_SEND_HOSTNAME,      OGRT_MSG_SEND_ENVIRONMENT);
    printf("  OGRT_MSG_SEND_CMDLINE=%d\n  OGRT_MSG_SEND_LOADEDMODULES=%d\n",
              OGRT_MSG_SEND_CMDLINE,      OGRT_MSG_SEND_LOADEDMODULES);
#if OGRT_FILTER_REGEXPS == 1
    printf("  OGRT_FILTER_REGEXP_LIST:\n");
    char *regexp_list[] = { OGRT_FILTER_REGEXPS_LIST };
    for(int i=0; i < OGRT_FILTER_REGEXPS_LIST_LENGTH; i++) {
      printf("    %s\n", regexp_list[i]);
    }
#endif
#if OGRT_MSG_SEND_ENVIRONMENT == 1 && defined(OGRT_MSG_SEND_ENVIRONMENT_WHITELIST)
    printf("  OGRT_ENVIRONMENT_WHITELIST:\n");
    char *whitelist[] = { OGRT_MSG_SEND_ENVIRONMENT_WHITELIST };
    for(int i=0; i < OGRT_MSG_SEND_ENVIRONMENT_WHITELIST_LENGTH; i++) {
      printf("    %s\n", whitelist[i]);
    }
#endif
  }

  if(ogrt_env_enabled("OGRT_ACTIVE")) {
    __ogrt_active = true;
  }

  if(__ogrt_active && __daemon_socket < 0) {
    /* cache PID of current process - we are reusing that quite often */
    __pid = getpid();
    __parent_pid = getppid();
    if(gethostname(__hostname, HOST_NAME_MAX) != 0) {
      Log(OGRT_LOG_ERR, "%s\n", strerror(errno));
      __ogrt_active = false;
      return 1;
    }

#if OGRT_FILTER_REGEXPS == 1
    char *regexps[] = { OGRT_FILTER_REGEXPS_LIST };
    regex_t regex;
    for(int i=0; i < OGRT_FILTER_REGEXPS_LIST_LENGTH; i++){
      if(regcomp(&regex, regexps[i], 0)) {
        Log(OGRT_LOG_ERR, "Can not compile regex: %s\n", regexps[i]);
        __ogrt_active = false;
        return 1;
      }

      if(regexec(&regex, ogrt_get_binpath(__pid), 0, NULL, 0) == 0){
        Log(OGRT_LOG_DBG, "Matching regex: %s\n", regexps[i]);
        __ogrt_active = false;
        return 1;
      }
    }
#endif

    /* establish a connection the the ogrt server */
    struct addrinfo hints, *servinfo, *p;
    memset(&hints, 0, sizeof hints);
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_DGRAM;

    int ret;
    if ((ret = getaddrinfo(OGRT_NET_HOST, OGRT_NET_PORT, &hints, &servinfo)) != 0) {
      Log(OGRT_LOG_ERR, "getaddrinfo: %s\n", gai_strerror(ret));
      __ogrt_active = false;
      return 1;
    }

    for(p = servinfo; p != NULL; p = p->ai_next) {
        if ((__daemon_socket = socket(p->ai_family, p->ai_socktype, p->ai_protocol)) == -1) {
            Log(OGRT_LOG_ERR, "%s\n", strerror(errno));
            __ogrt_active = false;
            continue;
        }

        if (connect(__daemon_socket, p->ai_addr, p->ai_addrlen) == -1) {
            close(__daemon_socket);
            __ogrt_active = false;
            Log(OGRT_LOG_ERR, "%s\n", strerror(errno));
            continue;
        }

        break;
    }

    if (p == NULL) {
        Log(OGRT_LOG_ERR, "%s\n", strerror(errno));
        __ogrt_active = false;
        return 1;
    }

    freeaddrinfo(servinfo);

    // hook signal functions to be able to wrap with our handler
    real_sigaction = dlsym(RTLD_NEXT, "sigaction");
    Log(OGRT_LOG_DBG, "sigaction() pointer: %p\n", real_sigaction);
    if(real_sigaction == NULL) {
      Log(OGRT_LOG_ERR, "Error in dlsym(): %s\n", dlerror());
      __ogrt_active = 0;
      return 1;
    }
    real_signal = dlsym(RTLD_NEXT, "signal");
    Log(OGRT_LOG_DBG, "signal() pointer: %p\n", real_signal);
    if(real_signal == NULL) {
      Log(OGRT_LOG_ERR, "Error in dlsym(): %s\n", dlerror());
      __ogrt_active = 0;
      return 1;
    }

    // wrap signal handlers
    int wrap_signals[] = {
      SIGHUP, SIGINT, SIGQUIT, SIGILL, SIGABRT, SIGFPE, SIGSEGV, SIGPIPE, SIGALRM, SIGTERM,
      SIGUSR1, SIGUSR2, SIGBUS, SIGPOLL, SIGPROF, SIGSYS, SIGTRAP, SIGVTALRM, SIGXCPU, SIGXFSZ,
      -1
    };
    for(int i=0; wrap_signals[i] != -1; i++) {
      struct sigaction old;
      int current_signal = wrap_signals[i];

      if((*real_sigaction)(current_signal, NULL, &old) == -1) {
        Log(OGRT_LOG_ERR, "got no signal handler");
        __ogrt_active = 0;
        return 1;
      }

      saved_signal_handlers[current_signal] = old.sa_handler;
      Log(OGRT_LOG_DBG, "%d: old signal handler is %p\n", current_signal, old.sa_handler);
      old.sa_handler = signal_wrapper;

      if((*real_sigaction)(current_signal, &old, NULL) == -1) {
        Log(OGRT_LOG_ERR, "could not set signal handler");
        __ogrt_active = 0;
        return 1;
      }
    }

    Log(OGRT_LOG_INFO, "I be watchin' yo! (process %d [%s] with parent %d)\n", __pid, ogrt_get_binpath(__pid), getppid());

    if(!ogrt_send_processinfo()) {
      Log(OGRT_LOG_ERR, "failed to send process info\n");
      __ogrt_active = 0;
      return 1;
    }
  }

  return 0;
}

ON_FINISH
int ogrt_preload_finish_hook() {
    int stderr_fp;
    if (dup2(stderr_fp, STDERR_FILENO) != -1) {
        ogrt_log_file = fdopen(stderr_fp, "w");
    }
    Log(OGRT_LOG_DBG, "called finish hook\n");

    if(__ogrt_active == 1) {
      if(!ogrt_send_resourceinfo()) {
            Log(OGRT_LOG_ERR, "failed to send resource info\n");
            return 1;
      }
    }
    return 0;
}

OGRT_INTERNAL
bool ogrt_send_resourceinfo() {
  Log(OGRT_LOG_DBG, "sending resource info...\n");

  OGRT__ProcessResourceInfo msg;
  ogrt__process_resource_info__init(&msg);

  struct timespec ts;
  clock_gettime(CLOCK_REALTIME, &ts);
  msg.time = (ts.tv_sec * (uint64_t)1000) + (ts.tv_nsec / 1000000);

  msg.uuid.data = __uuid;
  msg.uuid.len = 16;

  struct rusage ru;
  int ret = getrusage(RUSAGE_SELF, &ru);
  if (ret == 0) {
    msg.ru_utime    = (ru.ru_utime.tv_sec * (uint64_t)1000) + (ru.ru_utime.tv_usec / 1000);
    msg.ru_stime    = (ru.ru_stime.tv_sec * (uint64_t)1000) + (ru.ru_stime.tv_usec / 1000);
    msg.ru_maxrss   = ru.ru_maxrss;
    msg.ru_minflt   = ru.ru_minflt;
    msg.ru_majflt   = ru.ru_majflt;
    msg.ru_inblock  = ru.ru_inblock;
    msg.ru_oublock  = ru.ru_oublock;
    msg.ru_nvcsw    = ru.ru_nvcsw;
    msg.ru_nivcsw   = ru.ru_nivcsw;
  }

  memset(&ru, 0, sizeof(struct rusage));
  ret = getrusage(RUSAGE_CHILDREN, &ru);
  if (ret == 0) {
    msg.ru_utime    += (ru.ru_utime.tv_sec * (uint64_t)1000) + (ru.ru_utime.tv_usec / 1000);
    msg.ru_stime    += (ru.ru_stime.tv_sec * (uint64_t)1000) + (ru.ru_stime.tv_usec / 1000);
    msg.ru_maxrss   += ru.ru_maxrss;
    msg.ru_minflt   += ru.ru_minflt;
    msg.ru_majflt   += ru.ru_majflt;
    msg.ru_inblock  += ru.ru_inblock;
    msg.ru_oublock  += ru.ru_oublock;
    msg.ru_nvcsw    += ru.ru_nvcsw;
    msg.ru_nivcsw   += ru.ru_nivcsw;
  }


  size_t msg_len = ogrt__process_resource_info__get_packed_size(&msg);
  void *msg_serialized = NULL;
  char *msg_buffer = NULL;
  int send_length = ogrt_prepare_sendbuffer(OGRT__MESSAGE_TYPE__ProcessResourceMsg, msg_len, &msg_buffer, &msg_serialized);

  ogrt__process_resource_info__pack(&msg, msg_serialized);
  send(__daemon_socket, msg_buffer, send_length, 0);

  free(msg_buffer);
  return true;
}

OGRT_INTERNAL
bool ogrt_send_processinfo() {
    //TODO: refactor the process.
    // it is kind of dirty. the currently running binary is not an so.

    so_infos *so_infos = ogrt_get_loaded_so();
    OGRT__SharedObject *shared_object_excl_blank = &(so_infos->shared_objects[2]);

    //fprintf(stderr, "OGRT: Listing shared objects:\n");

    OGRT__SharedObject *shared_object_ptr[so_infos->size];
    //for(int i = 0; i < *so_info_size; i++) {
    //  ogrt_log_debug("[D] shared object path=%s, signature=%s\n", shared_object[i].path, shared_object[i].signature);
    //  fprintf(stderr, "OGRT:\tshared object path=%s, signature=%s\n", shared_object[i].path, shared_object[i].signature);
    //  shared_object_ptr[i] = &(shared_object[i]);
    //}
    for(int i = 0; i < so_infos->size-2; i++) {
      shared_object_ptr[i] = &(shared_object_excl_blank[i]);
    }

    struct timespec ts;
    clock_gettime(CLOCK_REALTIME, &ts);

    OGRT__ProcessInfo msg;
    ogrt__process_info__init(&msg);
    msg.binpath = ogrt_get_binpath(__pid);

#if OGRT_MSG_SEND_LOADEDMODULES == 1
    OGRT__Module **loaded_modules = NULL;
    char *module_env = getenv("LOADEDMODULES");
    msg.n_loaded_modules = 0;
    if(module_env) {
      /* first count the modules */
      int module_count = 0;
      char *count_string = strdup(module_env);
      char *count_token = strtok(count_string, ":");
      while(count_token) { module_count++; count_token = strtok(NULL, ":"); }
      free(count_string);

      Log(OGRT_LOG_DBG, "[D] checking env variable LOADEDMODULES with %d occurrences\n", module_count);

      if(module_count > 0) {
        /* allocate space for modules */
        loaded_modules = malloc(sizeof(OGRT__Module) * module_count);

        /* fill the protobuf */
        char *module_string = strdup(module_env);
        char *module_token = strtok(module_string, ":");
        for (module_count=0; module_token; module_count++) {
          Log(OGRT_LOG_DBG, "[D] token iteration %d\n", module_count);

          loaded_modules[module_count] = malloc(sizeof(OGRT__Module));
          ogrt__module__init(loaded_modules[module_count]);
          loaded_modules[module_count]->name = strdup(module_token);

          Log(OGRT_LOG_DBG, "[D] %s module detected\n", loaded_modules[module_count]->name);
          module_token=strtok(NULL, ":");
        }
        free(module_string);
        msg.n_loaded_modules = module_count;
        Log(OGRT_LOG_DBG, "[D] module count: %ld \n", msg.n_loaded_modules);
      }
    }

    if(msg.n_loaded_modules > 0) {
      msg.loaded_modules = loaded_modules;
    }
#endif

#if OGRT_MSG_SEND_CMDLINE == 1
    char *cmdline = ogrt_get_cmdline(__pid);
    msg.cmdline = cmdline;
#endif
    msg.pid = __pid;
    msg.parent_pid = __parent_pid;
    msg.time = (ts.tv_sec * (uint64_t)1000) + (ts.tv_nsec / 1000000);
    char *job_id = getenv(OGRT_ENV_JOBID);
    msg.job_id = job_id == NULL ? "UNKNOWN" : job_id;
#if OGRT_MSG_SEND_USERNAME == 1
    char *username = ogrt_get_username();
    msg.username = username == NULL ? "UNKNOWN" : username;
#endif
#if OGRT_MSG_SEND_HOSTNAME == 1
    char *hostname= ogrt_get_hostname();
    msg.hostname = hostname == NULL ? "UNKNOWN" : hostname;
#endif
#if OGRT_MSG_SEND_ENVIRONMENT == 1
    size_t envvar_count = 0;
    #ifdef OGRT_MSG_SEND_ENVIRONMENT_WHITELIST
    /* only get whitelisted variables */
    char *environment[OGRT_MSG_SEND_ENVIRONMENT_WHITELIST_LENGTH+1];
    char *whitelist[] = { OGRT_MSG_SEND_ENVIRONMENT_WHITELIST };

    for(int i=0; i < OGRT_MSG_SEND_ENVIRONMENT_WHITELIST_LENGTH; i++) {
      Log(OGRT_LOG_DBG, "[D] checking env variable: %s\n", whitelist[i]);
      char *env = getenv(whitelist[i]);
      if(env != NULL) {
        Log(OGRT_LOG_DBG, "[D] storing : %s with value '%s'\n", whitelist[i], env);
        int ret = asprintf(&(environment[envvar_count++]), "%s=%s", whitelist[i], env);
        if(ret == -1) {
          Log(OGRT_LOG_ERR, "failed copying environment variable\n");
        }
      }
    }
    environment[envvar_count] = NULL;
    #else
    /* get the whole environment */
    char **environment = environ;
    for(char **iterator = environment; *iterator != NULL; iterator++){
      envvar_count++;
    }
    #endif
    msg.n_environment_variables = envvar_count;
    msg.environment_variables = environment;
#endif
    if(so_infos->shared_objects[0].signature != NULL) {
      msg.signature = so_infos->shared_objects[0].signature;
    }
    msg.n_shared_objects = so_infos->size-2;
    msg.shared_objects = shared_object_ptr;

    uuid_generate(__uuid);
    msg.uuid.data = __uuid;
    msg.uuid.len= 16;

    size_t msg_len = ogrt__process_info__get_packed_size(&msg);
    void *msg_serialized = NULL;
    char *msg_buffer = NULL;
    int send_length = ogrt_prepare_sendbuffer(OGRT__MESSAGE_TYPE__ProcessInfoMsg, msg_len, &msg_buffer, &msg_serialized);

    ogrt__process_info__pack(&msg, msg_serialized);
    send(__daemon_socket, msg_buffer, send_length, 0);

    /* free stuff */
    free(msg.binpath);
    for(int i=0; i < so_infos->size; i++) {
      free(so_infos->shared_objects[i].path);
    }
    free(so_infos);
    free(msg_buffer);

#if OGRT_MSG_SEND_LOADEDMODULES == 1
    if(msg.n_loaded_modules > 0) {
      for(int i = 0; i < msg.n_loaded_modules; i++) {
        free(loaded_modules[i]->name);
        free(loaded_modules[i]);
      }
      free(msg.loaded_modules);
    }
#endif
#if OGRT_MSG_SEND_USERNAME == 1
    free(username);
#endif
#if OGRT_MSG_SEND_HOSTNAME == 1
    free(hostname);
#endif
#if OGRT_MSG_SEND_CMDLINE == 1
    free(cmdline);
#endif
#if OGRT_MSG_SEND_ENVIRONMENT == 1
#ifdef OGRT_MSG_SEND_ENVIRONMENT_WHITELIST
    for(int i=0; i < envvar_count; i++) {
      free(environment[i]);
    }
#endif
#endif
    return true;
}

/**
 * Prepare send buffer for shipping to daemon.
 * Takes a message type and the length of the payload and return the beginning of
 * the buffer, the beginning of the payload and the total size of the buffer.
 * Message format:
 *       32bit            32bit                  up to 32bit length
 * +----------------+----------------+--------------------------------------------+
 * |  message type  | payload_length |                payload                     |
 * +----------------+----------------+--------------------------------------------+
 *
 * This function is incredibly ugly. Should be reworked, but it works, right?
 */
OGRT_INTERNAL
int ogrt_prepare_sendbuffer(const int message_type, const int payload_length, char **buffer, void **payload) {
  uint32_t type = htonl(message_type);
  uint32_t length = htonl(payload_length);
  int total_length = payload_length + sizeof(type) + sizeof(length);

  *buffer = malloc(total_length);
  *payload = (((char *)*buffer) + sizeof(type) + sizeof(length));

  memcpy(*buffer, &type, sizeof(type));
  memcpy(*buffer + 4, &length, sizeof(length));

  return total_length;
}

OGRT_INTERNAL
void signal_wrapper(int signum) {
    int stderr_fp;
    if (dup2(stderr_fp, STDERR_FILENO) != -1) {
        ogrt_log_file = fdopen(stderr_fp, "w");
    }
    Log(OGRT_LOG_DBG, "in wrapper for signal %d\n", signum);

    if(__ogrt_active == 1) {
      if(!ogrt_send_resourceinfo()) {
        Log(OGRT_LOG_ERR, "failed to send resource info\n");
        return;
      }
    }

    if(saved_signal_handlers[signum] != NULL) {
        Log(OGRT_LOG_DBG, "calling saved handler for %d\n", signum);
        (*saved_signal_handlers[signum])(signum);
        return;
    }

    exit(128+signum);
}

/** Hook Functions **/
int sigaction(int signum, const struct sigaction *act, struct sigaction *oldact) {
  Log(OGRT_LOG_DBG, "hooked sigaction for signal %d\n", signum);

  if(act != NULL) {
    if(act->sa_handler == SIG_IGN) {
      Log(OGRT_LOG_DBG, "ignoring signal %d\n", signum);
      return real_sigaction(signum, act, oldact);
    }
    Log(OGRT_LOG_DBG, "%d: wrapping function %p\n", signum, act->sa_handler);
    saved_signal_handlers[signum] = act->sa_handler;
    ((struct sigaction *)act)->sa_handler = signal_wrapper;
  }
  return real_sigaction(signum, (const struct sigaction *)act, oldact);
}

sighandler_t signal (int signum, sighandler_t action) {
  Log(OGRT_LOG_DBG, "hooked signal() for signal %d\n", signum);
  if(action == SIG_IGN) {
    Log(OGRT_LOG_DBG, "ignoring signal %d\n", signum);
    return real_signal(signum, action);
  }

  Log(OGRT_LOG_DBG, "%d: wrapping function %p\n", action);
  saved_signal_handlers[signum] = action;
  action = signal_wrapper;
  return real_signal(signum, action);
}
