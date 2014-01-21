#ifndef LMCTFY_C_BINDING_LMCTFY_C_H_
#define LMCTFY_C_BINDING_LMCTFY_C_H_
#include "lmctfy.pb-c.h"
#include "status-c.h"

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

struct container;
struct container_api;

struct status *lmctfy_init_machine_raw(void *spec, int spec_size);
struct status *lmctfy_init_machine(Containers__Lmctfy__InitSpec *spec);
struct statuc *lmctfy_new_container_api(struct container_api **api);

#ifdef __cplusplus
}
#endif // __cplusplus
#endif // LMCTFY_C_BINDING_LMCTFY_C_H_