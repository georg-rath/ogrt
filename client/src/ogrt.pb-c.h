/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: protocol/ogrt.proto */

#ifndef PROTOBUF_C_protocol_2fogrt_2eproto__INCLUDED
#define PROTOBUF_C_protocol_2fogrt_2eproto__INCLUDED

#include <protobuf-c/protobuf-c.h>

PROTOBUF_C__BEGIN_DECLS

#if PROTOBUF_C_VERSION_NUMBER < 1003000
# error This file was generated by a newer version of protoc-c which is incompatible with your libprotobuf-c headers. Please update your headers.
#elif 1003001 < PROTOBUF_C_MIN_COMPILER_VERSION
# error This file was generated by an older version of protoc-c which is incompatible with your libprotobuf-c headers. Please regenerate this file with a newer version of protoc-c.
#endif


typedef struct _OGRT__JobStart OGRT__JobStart;
typedef struct _OGRT__JobEnd OGRT__JobEnd;
typedef struct _OGRT__SharedObject OGRT__SharedObject;
typedef struct _OGRT__Module OGRT__Module;
typedef struct _OGRT__ProcessInfo OGRT__ProcessInfo;
typedef struct _OGRT__ProcessResourceInfo OGRT__ProcessResourceInfo;
typedef struct _OGRT__JobInfo OGRT__JobInfo;


/* --- enums --- */

typedef enum _OGRT__MessageType {
  OGRT__MESSAGE_TYPE__JobStartMsg = 0,
  OGRT__MESSAGE_TYPE__JobEndMsg = 11,
  OGRT__MESSAGE_TYPE__ProcessInfoMsg = 12,
  OGRT__MESSAGE_TYPE__ProcessResourceMsg = 16,
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
    , (char *)protobuf_c_empty_string, 0 }


struct  _OGRT__JobEnd
{
  ProtobufCMessage base;
  char *job_id;
  int64_t end_time;
};
#define OGRT__JOB_END__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__job_end__descriptor) \
    , (char *)protobuf_c_empty_string, 0 }


struct  _OGRT__SharedObject
{
  ProtobufCMessage base;
  char *path;
  char *signature;
};
#define OGRT__SHARED_OBJECT__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__shared_object__descriptor) \
    , (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string }


struct  _OGRT__Module
{
  ProtobufCMessage base;
  char *name;
};
#define OGRT__MODULE__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__module__descriptor) \
    , (char *)protobuf_c_empty_string }


/*
 * sent at start of process 
 */
struct  _OGRT__ProcessInfo
{
  ProtobufCMessage base;
  ProtobufCBinaryData uuid;
  char *binpath;
  int32_t pid;
  int32_t parent_pid;
  int64_t time;
  char *signature;
  char *job_id;
  char *username;
  char *hostname;
  char *cmdline;
  char *cwd;
  size_t n_environment_variables;
  char **environment_variables;
  size_t n_arguments;
  char **arguments;
  size_t n_shared_objects;
  OGRT__SharedObject **shared_objects;
  size_t n_loaded_modules;
  OGRT__Module **loaded_modules;
};
#define OGRT__PROCESS_INFO__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__process_info__descriptor) \
    , {0,NULL}, (char *)protobuf_c_empty_string, 0, 0, 0, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, 0,NULL, 0,NULL, 0,NULL, 0,NULL }


/*
 * sent at end of process 
 */
struct  _OGRT__ProcessResourceInfo
{
  ProtobufCMessage base;
  ProtobufCBinaryData uuid;
  int64_t time;
  /*
   * resource info from getrusage() 
   */
  /*
   * user CPU time used 
   */
  int64_t ru_utime;
  /*
   * system CPU time used 
   */
  int64_t ru_stime;
  /*
   * maximum resident set size 
   */
  int64_t ru_maxrss;
  /*
   * page reclaims (soft page faults) 
   */
  int64_t ru_minflt;
  /*
   * page faults (hard page faults) 
   */
  int64_t ru_majflt;
  /*
   * block input operations 
   */
  int64_t ru_inblock;
  /*
   * block output operations 
   */
  int64_t ru_oublock;
  /*
   * voluntary context switches 
   */
  int64_t ru_nvcsw;
  /*
   * involuntary context switches 
   */
  int64_t ru_nivcsw;
};
#define OGRT__PROCESS_RESOURCE_INFO__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__process_resource_info__descriptor) \
    , {0,NULL}, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }


struct  _OGRT__JobInfo
{
  ProtobufCMessage base;
  char *job_id;
  size_t n_processes;
  OGRT__ProcessInfo **processes;
};
#define OGRT__JOB_INFO__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&ogrt__job_info__descriptor) \
    , (char *)protobuf_c_empty_string, 0,NULL }


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
/* OGRT__Module methods */
void   ogrt__module__init
                     (OGRT__Module         *message);
size_t ogrt__module__get_packed_size
                     (const OGRT__Module   *message);
size_t ogrt__module__pack
                     (const OGRT__Module   *message,
                      uint8_t             *out);
size_t ogrt__module__pack_to_buffer
                     (const OGRT__Module   *message,
                      ProtobufCBuffer     *buffer);
OGRT__Module *
       ogrt__module__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__module__free_unpacked
                     (OGRT__Module *message,
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
/* OGRT__ProcessResourceInfo methods */
void   ogrt__process_resource_info__init
                     (OGRT__ProcessResourceInfo         *message);
size_t ogrt__process_resource_info__get_packed_size
                     (const OGRT__ProcessResourceInfo   *message);
size_t ogrt__process_resource_info__pack
                     (const OGRT__ProcessResourceInfo   *message,
                      uint8_t             *out);
size_t ogrt__process_resource_info__pack_to_buffer
                     (const OGRT__ProcessResourceInfo   *message,
                      ProtobufCBuffer     *buffer);
OGRT__ProcessResourceInfo *
       ogrt__process_resource_info__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   ogrt__process_resource_info__free_unpacked
                     (OGRT__ProcessResourceInfo *message,
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
typedef void (*OGRT__Module_Closure)
                 (const OGRT__Module *message,
                  void *closure_data);
typedef void (*OGRT__ProcessInfo_Closure)
                 (const OGRT__ProcessInfo *message,
                  void *closure_data);
typedef void (*OGRT__ProcessResourceInfo_Closure)
                 (const OGRT__ProcessResourceInfo *message,
                  void *closure_data);
typedef void (*OGRT__JobInfo_Closure)
                 (const OGRT__JobInfo *message,
                  void *closure_data);

/* --- services --- */


/* --- descriptors --- */

extern const ProtobufCEnumDescriptor    ogrt__message_type__descriptor;
extern const ProtobufCMessageDescriptor ogrt__job_start__descriptor;
extern const ProtobufCMessageDescriptor ogrt__job_end__descriptor;
extern const ProtobufCMessageDescriptor ogrt__shared_object__descriptor;
extern const ProtobufCMessageDescriptor ogrt__module__descriptor;
extern const ProtobufCMessageDescriptor ogrt__process_info__descriptor;
extern const ProtobufCMessageDescriptor ogrt__process_resource_info__descriptor;
extern const ProtobufCMessageDescriptor ogrt__job_info__descriptor;

PROTOBUF_C__END_DECLS


#endif  /* PROTOBUF_C_protocol_2fogrt_2eproto__INCLUDED */
