/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: protocol/ogrt.proto */

#ifndef PROTOBUF_C_protocol_2fogrt_2eproto__INCLUDED
#define PROTOBUF_C_protocol_2fogrt_2eproto__INCLUDED

#include <protobuf-c/protobuf-c.h>

PROTOBUF_C__BEGIN_DECLS

#if PROTOBUF_C_VERSION_NUMBER < 1000000
# error This file was generated by a newer version of protoc-c which is incompatible with your libprotobuf-c headers. Please update your headers.
#elif 1001001 < PROTOBUF_C_MIN_COMPILER_VERSION
# error This file was generated by an older version of protoc-c which is incompatible with your libprotobuf-c headers. Please regenerate this file with a newer version of protoc-c.
#endif


typedef struct _OGRT__JobStart OGRT__JobStart;
typedef struct _OGRT__JobEnd OGRT__JobEnd;
typedef struct _OGRT__SharedObject OGRT__SharedObject;
typedef struct _OGRT__ProcessInfo OGRT__ProcessInfo;
typedef struct _OGRT__JobInfo OGRT__JobInfo;
typedef struct _OGRT__Fork OGRT__Fork;
typedef struct _OGRT__Execve OGRT__Execve;


/* --- enums --- */

typedef enum _OGRT__MessageType {
  OGRT__MESSAGE_TYPE__JobStartMsg = 10,
  OGRT__MESSAGE_TYPE__JobEndMsg = 11,
  OGRT__MESSAGE_TYPE__ProcessInfoMsg = 12,
  OGRT__MESSAGE_TYPE__SharedObjectMsg = 13,
  OGRT__MESSAGE_TYPE__ForkMsg = 14,
  OGRT__MESSAGE_TYPE__ExecveMsg = 15
    PROTOBUF_C__FORCE_ENUM_TO_BE_INT_SIZE(OGRT__MESSAGE_TYPE)
} OGRT__MessageType;

/* --- messages --- */

struct  _OGRT__JobStart
{
  ProtobufCMessage base;
  char *job_id;
  int64_t start_time;
};
#define OGRT__JOB_START__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__job_start__descriptor) \
    , NULL, 0 }


struct  _OGRT__JobEnd
{
  ProtobufCMessage base;
  char *job_id;
  int64_t end_time;
};
#define OGRT__JOB_END__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__job_end__descriptor) \
    , NULL, 0 }


struct  _OGRT__SharedObject
{
  ProtobufCMessage base;
  char *path;
  char *signature;
};
#define OGRT__SHARED_OBJECT__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__shared_object__descriptor) \
    , NULL, NULL }


struct  _OGRT__ProcessInfo
{
  ProtobufCMessage base;
  char *binpath;
  int32_t pid;
  int32_t parent_pid;
  int64_t time;
  char *signature;
  char *job_id;
  char *username;
  char *hostname;
  char *cmdline;
  size_t n_environment_variables;
  char **environment_variables;
  size_t n_arguments;
  char **arguments;
  size_t n_shared_objects;
  OGRT__SharedObject **shared_objects;
};
#define OGRT__PROCESS_INFO__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__process_info__descriptor) \
    , NULL, 0, 0, 0, NULL, NULL, NULL, NULL, NULL, 0,NULL, 0,NULL, 0,NULL }


struct  _OGRT__JobInfo
{
  ProtobufCMessage base;
  char *job_id;
  size_t n_processes;
  OGRT__ProcessInfo **processes;
};
#define OGRT__JOB_INFO__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__job_info__descriptor) \
    , NULL, 0,NULL }


struct  _OGRT__Fork
{
  ProtobufCMessage base;
  char *hostname;
  int32_t parent_pid;
  int32_t child_pid;
  char *name;
};
#define OGRT__FORK__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__fork__descriptor) \
    , NULL, 0, 0, NULL }


struct  _OGRT__Execve
{
  ProtobufCMessage base;
  char *hostname;
  int32_t pid;
  int32_t parent_pid;
  char *filename;
  size_t n_arguments;
  char **arguments;
  size_t n_environment_variables;
  char **environment_variables;
  char *uuid;
};
#define OGRT__EXECVE__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__execve__descriptor) \
    , NULL, 0, 0, NULL, 0,NULL, 0,NULL, NULL }


/* OGRT__JobStart methods */
void   ogrt__job_start__init
                     (OGRT__JobStart         *message);
size_t ogrt__job_start__get_packed_size
                     (const OGRT__JobStart   *message);
size_t ogrt__job_start__pack
                     (const OGRT__JobStart   *message,
                      uint8_t             *out);
size_t ogrt__job_start__pack_to_buffer
                     (const OGRT__JobStart   *message,
                      ProtobufCBuffer     *buffer);
OGRT__JobStart *
       ogrt__job_start__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__job_start__free_unpacked
                     (OGRT__JobStart *message,
                      ProtobufCAllocator *allocator);
/* OGRT__JobEnd methods */
void   ogrt__job_end__init
                     (OGRT__JobEnd         *message);
size_t ogrt__job_end__get_packed_size
                     (const OGRT__JobEnd   *message);
size_t ogrt__job_end__pack
                     (const OGRT__JobEnd   *message,
                      uint8_t             *out);
size_t ogrt__job_end__pack_to_buffer
                     (const OGRT__JobEnd   *message,
                      ProtobufCBuffer     *buffer);
OGRT__JobEnd *
       ogrt__job_end__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__job_end__free_unpacked
                     (OGRT__JobEnd *message,
                      ProtobufCAllocator *allocator);
/* OGRT__SharedObject methods */
void   ogrt__shared_object__init
                     (OGRT__SharedObject         *message);
size_t ogrt__shared_object__get_packed_size
                     (const OGRT__SharedObject   *message);
size_t ogrt__shared_object__pack
                     (const OGRT__SharedObject   *message,
                      uint8_t             *out);
size_t ogrt__shared_object__pack_to_buffer
                     (const OGRT__SharedObject   *message,
                      ProtobufCBuffer     *buffer);
OGRT__SharedObject *
       ogrt__shared_object__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__shared_object__free_unpacked
                     (OGRT__SharedObject *message,
                      ProtobufCAllocator *allocator);
/* OGRT__ProcessInfo methods */
void   ogrt__process_info__init
                     (OGRT__ProcessInfo         *message);
size_t ogrt__process_info__get_packed_size
                     (const OGRT__ProcessInfo   *message);
size_t ogrt__process_info__pack
                     (const OGRT__ProcessInfo   *message,
                      uint8_t             *out);
size_t ogrt__process_info__pack_to_buffer
                     (const OGRT__ProcessInfo   *message,
                      ProtobufCBuffer     *buffer);
OGRT__ProcessInfo *
       ogrt__process_info__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__process_info__free_unpacked
                     (OGRT__ProcessInfo *message,
                      ProtobufCAllocator *allocator);
/* OGRT__JobInfo methods */
void   ogrt__job_info__init
                     (OGRT__JobInfo         *message);
size_t ogrt__job_info__get_packed_size
                     (const OGRT__JobInfo   *message);
size_t ogrt__job_info__pack
                     (const OGRT__JobInfo   *message,
                      uint8_t             *out);
size_t ogrt__job_info__pack_to_buffer
                     (const OGRT__JobInfo   *message,
                      ProtobufCBuffer     *buffer);
OGRT__JobInfo *
       ogrt__job_info__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__job_info__free_unpacked
                     (OGRT__JobInfo *message,
                      ProtobufCAllocator *allocator);
/* OGRT__Fork methods */
void   ogrt__fork__init
                     (OGRT__Fork         *message);
size_t ogrt__fork__get_packed_size
                     (const OGRT__Fork   *message);
size_t ogrt__fork__pack
                     (const OGRT__Fork   *message,
                      uint8_t             *out);
size_t ogrt__fork__pack_to_buffer
                     (const OGRT__Fork   *message,
                      ProtobufCBuffer     *buffer);
OGRT__Fork *
       ogrt__fork__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__fork__free_unpacked
                     (OGRT__Fork *message,
                      ProtobufCAllocator *allocator);
/* OGRT__Execve methods */
void   ogrt__execve__init
                     (OGRT__Execve         *message);
size_t ogrt__execve__get_packed_size
                     (const OGRT__Execve   *message);
size_t ogrt__execve__pack
                     (const OGRT__Execve   *message,
                      uint8_t             *out);
size_t ogrt__execve__pack_to_buffer
                     (const OGRT__Execve   *message,
                      ProtobufCBuffer     *buffer);
OGRT__Execve *
       ogrt__execve__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__execve__free_unpacked
                     (OGRT__Execve *message,
                      ProtobufCAllocator *allocator);
/* --- per-message closures --- */

typedef void (*OGRT__JobStart_Closure)
                 (const OGRT__JobStart *message,
                  void *closure_data);
typedef void (*OGRT__JobEnd_Closure)
                 (const OGRT__JobEnd *message,
                  void *closure_data);
typedef void (*OGRT__SharedObject_Closure)
                 (const OGRT__SharedObject *message,
                  void *closure_data);
typedef void (*OGRT__ProcessInfo_Closure)
                 (const OGRT__ProcessInfo *message,
                  void *closure_data);
typedef void (*OGRT__JobInfo_Closure)
                 (const OGRT__JobInfo *message,
                  void *closure_data);
typedef void (*OGRT__Fork_Closure)
                 (const OGRT__Fork *message,
                  void *closure_data);
typedef void (*OGRT__Execve_Closure)
                 (const OGRT__Execve *message,
                  void *closure_data);

/* --- services --- */


/* --- descriptors --- */

extern const ProtobufCEnumDescriptor    ogrt__message_type__descriptor;
extern const ProtobufCMessageDescriptor ogrt__job_start__descriptor;
extern const ProtobufCMessageDescriptor ogrt__job_end__descriptor;
extern const ProtobufCMessageDescriptor ogrt__shared_object__descriptor;
extern const ProtobufCMessageDescriptor ogrt__process_info__descriptor;
extern const ProtobufCMessageDescriptor ogrt__job_info__descriptor;
extern const ProtobufCMessageDescriptor ogrt__fork__descriptor;
extern const ProtobufCMessageDescriptor ogrt__execve__descriptor;

PROTOBUF_C__END_DECLS


#endif  /* PROTOBUF_C_protocol_2fogrt_2eproto__INCLUDED */
