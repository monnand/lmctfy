// Mock functions for clmctfy
#include "clmctfy.h"
#include "clmctfy-raw.h"
#include "macros.h"
#include <stdlib.h>

#define CONTAINER_API_ADDR	(struct container_api *)0x1234

struct container_api {
};

struct container {
	char *name;
};

static notification_id_t next_notif_id = 0;

MOCK_FUNCTION_BEGIN(lmctfy_init_machine_raw, const void *spec, const size_t spec_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_new_container_api, struct container_api **api) {
	// Some random pointer.
	// Will not be dereferenced.
	*api = CONTAINER_API_ADDR;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_api_create_container_raw,
    struct container **container,
    struct container_api *api,
    const char *container_name,
    const void *spec,
    const size_t spec_size) {
	*container = (struct container *)malloc(sizeof(struct container));
	(*container)->name = strdup(container_name);
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_run_raw,
                    pid_t *tid,
                    struct container *container,
                    const int argc,
                    const char **argv,
                    const void *spec,
                    const size_t spec_size) {
	*tid = 10;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_update_raw,
                    struct container *container,
                    int policy,
                    const void *spec,
                    const size_t spec_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_register_notification_raw,
                    notification_id_t *notif_id,
                    struct container *container,
                    lmctfy_event_callback_f callback,
                    void *user_data,
                    const void *spec,
                    const size_t spec_size) {
	*notif_id = next_notif_id++;
} MOCK_FUNCTION_END

void lmctfy_delete_container_api(struct container_api *api) {
	return;
}

void lmctfy_delete_container(struct container *container) {
	if (container != NULL) {
		if (container->name != NULL) {
			free(container->name);
		}
		free(container);
	}
}

MOCK_FUNCTION_BEGIN(lmctfy_container_api_get_container,
                    struct container **container,
                    const struct container_api *api,
                    const char *container_name) {
	*container = (struct container *)malloc(sizeof(struct container));
	(*container)->name = strdup(container_name);
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_api_destroy_container,
                    struct container_api *api,
                    struct container *container) {
} MOCK_FUNCTION_END


MOCK_FUNCTION_BEGIN(lmctfy_container_api_detect_container,
                    char **container_name,
                    struct container_api *api,
                    pid_t pid) {
	*container_name = strdup("/");
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_enter,
                    struct container *container,
                    const pid_t *tids,
                    const int n) {
} MOCK_FUNCTION_END

