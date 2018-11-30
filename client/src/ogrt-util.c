#include "ogrt-util.h"

/**
 * Check if an environment variable is true
 */
OGRT_INTERNAL
bool ogrt_env_enabled(char *env_variable) {
  char *env_var = getenv(env_variable);
  if(env_var != NULL && (strcmp(env_var, "yes") == 0 || strcmp(env_var, "true") == 0 || strcmp(env_var, "1") == 0)) {
    return true;
  }
  return false;
}

/**
 * Normalize a path to not contain '..' or '.'.
 * TODO: This function uses the glibc realpath internally and could do disk access.
 */
OGRT_INTERNAL
char *ogrt_normalize_path(const char *path) {
  char *normalized_path = malloc(PATH_MAX);
  if (normalized_path == NULL) {
    Log(OGRT_LOG_ERR, "failed  to allocate memory\n");
    return NULL;
  }
  char *ret = realpath(path, normalized_path);
  if(ret == NULL) {
    free(normalized_path);
    return NULL;
  }
  return normalized_path;
}

/**
 * Given the pid of a program return the path to its binary.
 * This is done by walking /proc.
 */
OGRT_INTERNAL
char *ogrt_get_binpath(const pid_t pid) {
  char proc_path[PATH_MAX];
  sprintf(proc_path, "/proc/%d/exe", pid);

  char *bin_path;
  bin_path = malloc(PATH_MAX);
  if (bin_path == NULL) {
    Log(OGRT_LOG_ERR, "failed  to allocate memory\n");
    return NULL;
  }

  ssize_t len = readlink(proc_path, bin_path, PATH_MAX);
  if (len == -1) {
     Log(OGRT_LOG_ERR, "failed getting executable path: %s\n", strerror(errno));
     free(bin_path);
     return NULL;
  }

  bin_path[len] = '\0';
  return bin_path;
}

/**
 * Given the pid of a program return the cmdline it
 * was called with.
 * This is done by walking /proc.
 */
OGRT_INTERNAL
char *ogrt_get_cmdline(const pid_t pid) {
  char proc_path[PATH_MAX];
  sprintf(proc_path, "/proc/%d/cmdline", pid);

  char *cmdline;
  cmdline = malloc(PATH_MAX);
  if (cmdline == NULL) {
    fprintf(stderr, "OGRT: memory allocate failed\n");
    return NULL;
  }

  int fd = open(proc_path, O_RDONLY);
  if(fd < 0) {
     Log(OGRT_LOG_ERR, "failed to get cmdline - open(): %s\n", strerror(errno));
     return NULL;
  }
  int nbytesread = read(fd, cmdline, PATH_MAX);
  if(nbytesread < 0) {
     Log(OGRT_LOG_ERR, "failed to get cmdline - read(): %s\n", strerror(errno));
     return NULL;
  }
  char *end = cmdline + nbytesread;

  /* cmdline is a string w/ NULLs instead of spaces, so
   * we replace them with spaces, except the last one */
  for(char *p = cmdline; p < (end-1); *p++) {
    if(*p == '\0') {
      *p = ' ';
    }
  }
  *end = '\0';

  return cmdline;
}

/**
 * Get username of the current user.
 * Does not use environment variables to do the lookup.
 */
OGRT_INTERNAL
char *ogrt_get_username() {
  struct passwd *pw = getpwuid(geteuid());
  if(pw == NULL) {
    Log(OGRT_LOG_ERR, "failed to get username: %s\n", strerror(errno));
    return NULL;
  }
  return strdup(pw->pw_name);
}

/**
 * Get the name of the current host.
 */
OGRT_INTERNAL
char *ogrt_get_hostname() {
  char hostname[HOST_NAME_MAX+1];
  int ret = gethostname(hostname, sizeof(hostname));
  if(ret) {
    Log(OGRT_LOG_ERR, "failed to get hostname - gethostname(): %s\n", strerror(errno));
    return NULL;
  }
  return strdup(hostname);
}


/**
 * Get current working directory
 */
OGRT_INTERNAL
char *ogrt_get_cwd() {
  char path[PATH_MAX];
  char *ret = getcwd(path, PATH_MAX);
  if(ret == NULL) {
    Log(OGRT_LOG_ERR, "failed to get cwd - getcwd(): %s\n", strerror(errno));
    return NULL;
  }
  return strdup(path);
}
