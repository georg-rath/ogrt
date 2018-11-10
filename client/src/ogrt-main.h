#ifndef OGRT_H_INCLUDED
#define OGRT_H_INCLUDED

#define _GNU_SOURCE
#include <dlfcn.h>
#include <sys/time.h>
#include <sys/resource.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <time.h>
#include <stdbool.h>
#include <uuid/uuid.h>
#include <errno.h>
#include <regex.h>
#include <signal.h>
#include "ogrt.pb-c.h"
#include "ogrt-readso.h"
#include "ogrt-log.h"
#include "ogrt-util.h"
#include "ogrt-cmdline.h"
#include "ogrt-signature.h"
#include "config.h"

/** ELF signature configuration */
#define OGRT_ELF_SECTION_NAME (".note.ogrt.info")
#define OGRT_ELF_NOTE_TYPE    (0x4f475254)

#define ON_INIT __attribute__((constructor))
#define ON_FINISH __attribute__((destructor))
#define OGRT_INTERNAL __attribute__ ((visibility ("hidden")))

/** function prototypes */
int ogrt_prepare_sendbuffer(const int message_type, const int message_length, char **buffer_begin, void **payload);
bool ogrt_send_processinfo();
bool ogrt_send_resourceinfo();
int ogrt_preload_init_hook();
void signal_wrapper(int signum);
#endif
