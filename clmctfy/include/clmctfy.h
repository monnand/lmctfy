#ifndef LMCTFY_C_BINDING_LMCTFY_C_H_
#define LMCTFY_C_BINDING_LMCTFY_C_H_
#include "lmctfy.pb-c.h"

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

struct status {
  int error_code;

  // Null-terminated string allocated on heap.
  // Needs to be free()'ed.
  char *message;
};

struct container;
struct container_api;

// Initializes the machine to start being able to create containers.
//
// Arguments:
//  - s: [output] s will be used as output. It contains the error code/message.
//  - spec: The specification.
//
// Returns:
//  Returns the error code. 0 on success. The return code is same as
//  status_get_code(s).
int lmctfy_init_machine(struct status *s, const Containers__Lmctfy__InitSpec *spec);

// Create a new container_api.
//
// Arguments:
//  - s: [output] s will be used as output. It contains the error code/message.
//  - api: [output] The address of a pointer to struct container_api. The
//  pointer of the container api will be stored in this address.
//
// Returns:
//  Returns the error code. 0 on success. The return code is same as
//  status_get_code(s).
int lmctfy_new_container_api(struct status *s, struct container_api **api);

// Release the container api.
void lmctfy_delete_container_api(struct container_api *api);

// Get a container
//
// Arguments:
//
//  - s: [output] s will be used as output. It contains the error code/message.
//  - container: [output] The address of a pointer to struct container. It will
//  be used to store the pointer to the container.
//  - api: A container api.
//  - container_name: the container name.
//
// Returns:
//  Returns the error code. 0 on success. The return code is same as
//  status_get_code(s).
int lmctfy_container_api_get_container(
    struct status *s,
    struct container **container,
    const struct container_api *api,
    const char *container_name);

#ifdef __cplusplus
}
#endif // __cplusplus
#endif // LMCTFY_C_BINDING_LMCTFY_C_H_