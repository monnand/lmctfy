// Mock functions for clmctfy
#include <stdlib.h>
#include "clmctfy.h"
#include "clmctfy-raw.h"
#include "clmctfy_mock.h"
#include "macros.h"

#define CONTAINER_API_ADDR	(struct container_api *)0x1234

struct container_api {
};

struct container {
	char *name;
};

static notification_id_t next_notif_id = 0;
void lmctfy_delete_container_api(struct container_api *api) {
	return;
}

const char *lmctfy_container_name(struct container *container) {
	return container->name;
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
                    const struct container_api *api,
                    const char *container_name,
                    struct container **container) {
	*container = (struct container *)malloc(sizeof(struct container));
	(*container)->name = strdup(container_name);
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_api_destroy_container,
                    struct container_api *api,
                    struct container *container) {
} MOCK_FUNCTION_END


MOCK_FUNCTION_BEGIN(lmctfy_container_api_detect_container,
                    struct container_api *api,
                    pid_t pid,
                    char **container_name) {
	*container_name = strdup("/");
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_enter,
                    struct container *container,
                    const pid_t *tids,
                    const int n) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_exec,
                    struct container *container,
                    const int argc,
                    const char **argv) {
} MOCK_FUNCTION_END

/*
MOCK_FUNCTION_BEGIN(lmctfy_container_spec,
                    struct container *container,
                    Containers__Lmctfy__ContainerSpec **spec) {
} MOCK_FUNCTION_END
*/

MOCK_FUNCTION_BEGIN(lmctfy_container_list_subcontainers,
                    struct container *container,
                    int list_policy,
                    struct container **subcontainers[],
                    int *subcontainers_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_list_threads,
                    struct container *container,
                    int list_policy,
                    pid_t *threads[],
                    int *threads_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_list_processes,
                    struct container *container,
                    int list_policy,
                    pid_t *processes[],
                    int *processes_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_pause,
                    struct container *container) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_resume,
                    struct container *container) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_killall,
                    struct container *container) {
} MOCK_FUNCTION_END

