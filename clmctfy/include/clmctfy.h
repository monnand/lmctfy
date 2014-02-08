#ifndef LMCTFY_C_BINDING_LMCTFY_C_H_
#define LMCTFY_C_BINDING_LMCTFY_C_H_

#include <unistd.h>
#include "lmctfy.pb-c.h"
#include "util/task/codes.pb-c.h"

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

enum {
  CONTAINER_UPDATE_POLICY_DIFF,
  CONTAINER_UPDATE_POLICY_REPLACE
};

enum {
  CONTAINER_LIST_POLICY_SELF,
  CONTAINER_LIST_POLICY_RECURSIVE
};

enum {
  CONTAINER_STATS_TYPE_SUMMARY,
  CONTAINER_STATS_TYPE_FULL
};

struct status {
  int error_code;

  // Null-terminated string allocated on heap.
  // Needs to be free()'ed.
  char *message;
};


struct container;
struct container_api;

// Callback used on an event notification.
//
// The container and status structure pointers are only valid within the
// callback.
//
// - container: The container that received the notification. It is an error
//    to delete it.
// - status: The status of the notification. If OK, then the event registered
//    occured. Otherwise, an error is reported in the status. Errors may
//    be caused by container deletion or unexpected registration errors.
typedef void lmctfy_event_callback_f(struct container *, const struct status *);

// Initializes the machine to start being able to create containers.
//
// Arguments:
//  - s: [output] s will be used as output.
//    If the original value of s->code is not
//    UTIL__ERROR__CODE__INVALID_ARGUMENT, then the function will return
//    immediately with value s->code. This will help users to propagete errors.
//    Otherwise, it contains the error code and message.
//  - spec: The specification.
//
// Returns:
//
//  Returns the error code. 0 on success. When there's an error, the return code
//  is same as s->error_code when s is not NULL.
int lmctfy_init_machine(struct status *s, const Containers__Lmctfy__InitSpec *spec);

// Create a new container_api.
//
// Arguments:
//  - s: [output] s will be used as output.
//    If the original value of s->code is not
//    UTIL__ERROR__CODE__INVALID_ARGUMENT, then the function will return
//    immediately with value s->code. This will help users to propagete errors.
//    Otherwise, it contains the error code and message.
//  - api: [output] The address of a pointer to struct container_api. The
//    pointer of the container api will be stored in this address.
//
// Returns:
//
//  Returns the error code. 0 on success. When there's an error, the return code
//  is same as s->error_code when s is not NULL.
int lmctfy_new_container_api(struct status *s, struct container_api **api);

// Release the container api. 
//
// Arguments:
//
//  - api: The container api. This pointer will be invalid after the call to
//  this function
void lmctfy_delete_container_api(struct container_api *api);

// Get a container
//
// Arguments:
//
//  - s: [output] s will be used as output.
//    If the original value of s->code is not
//    UTIL__ERROR__CODE__INVALID_ARGUMENT, then the function will return
//    immediately with value s->code. This will help users to propagete errors.
//    Otherwise, it contains the error code and message.
//  - container: [output] The address of a pointer to struct container. It will
//    be used to store the pointer to the container.
//  - api: A container api.
//  - container_name: the container name.
//
// Returns:
//
//  Returns the error code. 0 on success. When there's an error, the return code
//  is same as s->error_code when s is not NULL.
int lmctfy_container_api_get_container(
    struct status *s,
    struct container **container,
    const struct container_api *api,
    const char *container_name);

// Create a container
//
// Arguments:
//
//  - s: [output] s will be used as output.
//    If the original value of s->code is not
//    UTIL__ERROR__CODE__INVALID_ARGUMENT, then the function will return
//    immediately with value s->code. This will help users to propagete errors.
//    Otherwise, it contains the error code and message.
//  - container: [output] The address of a pointer to struct container. It will
//    be used to store the newly created container.
//  - api: A container api.
//  - container_name: the container name.
//  - spec: container specification.
//
// Returns:
//
//  Returns the error code. 0 on success. When there's an error, the return code
//  is same as s->error_code when s is not NULL.
int lmctfy_container_api_create_container(
    struct status *s,
    struct container **container,
    struct container_api *api,
    const char *container_name,
    const Containers__Lmctfy__ContainerSpec *spec);

// Destroy a container
//
// Arguments:
//
//  - s: [output] s will be used as output.
//    If the original value of s->code is not
//    UTIL__ERROR__CODE__INVALID_ARGUMENT, then the function will return
//    immediately with value s->code. This will help users to propagete errors.
//    Otherwise, it contains the error code and message.
//  - api: A container api.
//  - container: The pointer to struct container. The pointer will become
//    invalid after a success destroy().
//
// Returns:
//
//  Returns the error code. 0 on success. When there's an error, the return code
//  is same as s->error_code when s is not NULL.
int lmctfy_container_api_destroy_container(struct status *s,
                                           struct container_api *api,
                                           struct container *container);

int lmctfy_container_api_detect_container(struct status *s,
                                          char *container_name,
                                          const size_t n,
                                          struct container_api *api,
                                          pid_t pid);

// Release the memory used by the container structure. The refered container
// will not be affected.
//
// Arguments:
//
//  - container: The container.  This pointer will be invalid after the call to
//    this function
void lmctfy_delete_container(struct container *container);

int lmctfy_container_enter(struct status *s,
                           struct container *container,
                           const pid_t *tids,
                           const int n);

// 
// Arguments:
//
//  - s: [output]
//  - tid: [output] On success, tid stores the PID of the command.
//  - container
//  - argc: number of arguments (including the binary file path).
//  - argv: All arguments. The first element is the binary that will be executed
//    and must be an absolute path.
int lmctfy_container_run(struct status *s,
                         pid_t *tid,
                         struct container *container,
                         const int argc,
                         const char **argv,
                         const Containers__Lmctfy__RunSpec *spec);

int lmctfy_container_exec(struct status *s,
                          struct container *container,
                          const int argc,
                          const char **argv);

int lmctfy_container_update(struct status *s,
                            struct container *container,
                            int policy,
                            const Containers__Lmctfy__ContainerSpec *spec);

int lmctfy_container_spec(struct status *s,
                          struct container *container,
                          Containers__Lmctfy__ContainerSpec **spec);

// List all subcontainers.
//
// Arguments:
//
//  - s: [output]
//  - subcontainers: [output] The address of a pointer points to an array of
//    struct container pointers. The caller takes the onwership. On success, the
//    pointer will be assigned an address to an array of containers. The pointed 
//    array should be released with free(). The containers inside the array
//    should be deleted/destroyed individually.
//  - nr_subcontainers: [output] The address of an integer used to store number
//    of subcontainers in the array.
//  - container: The container
//  - list_policy: CONTAINER_LIST_POLICY_SELF or CONTAINER_LIST_POLICY_RECURSIVE
int lmctfy_container_list_subcontainers(struct status *s,
                                        struct container **subcontainers[],
                                        int *nr_subcontainers,
                                        struct container *container,
                                        int list_policy);
#ifdef __cplusplus
}
#endif // __cplusplus
#endif // LMCTFY_C_BINDING_LMCTFY_C_H_
