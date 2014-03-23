// Mock functions for clmctfy
#include "clmctfy.h"
#include <stdlib.h>
#include "clmctfy-raw.h"
#include "clmctfy_mock.h"
#include "macros.h"

#define CONTAINER_API_ADDR	(struct container_api *)0x1234

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
	(*container)->callback = NULL;
	(*container)->cb_userdata = NULL;
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


MOCK_FUNCTION_BEGIN(lmctfy_container_list_subcontainers,
                    struct container *container,
                    int list_policy,
                    struct container **subcontainers[],
                    int *subcontainers_size) {
	struct container *c;
	*subcontainers_size = 2;
	*subcontainers = (struct container **)malloc(sizeof(struct container *) * (*subcontainers_size));

	c = (struct container *)malloc(sizeof(struct container));
	c->name = strdup("/a");
	c->callback = NULL;
	c->cb_userdata = NULL;
	(*subcontainers)[0] = c;

	c = (struct container *)malloc(sizeof(struct container));
	c->name = strdup("/b");
	c->callback = NULL;
	c->cb_userdata = NULL;
	(*subcontainers)[1] = c;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_list_threads,
                    struct container *container,
                    int list_policy,
                    pid_t *threads[],
                    int *threads_size) {
	*threads_size = 2;
	*threads = (pid_t *)malloc(sizeof(pid_t) * (*threads_size));
	(*threads)[0] = 1;
	(*threads)[1] = 2;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_list_processes,
                    struct container *container,
                    int list_policy,
                    pid_t *processes[],
                    int *processes_size) {
	*processes_size = 3;
	*processes = (pid_t *)malloc(sizeof(pid_t) * (*processes_size));
	(*processes)[0] = 3;
	(*processes)[1] = 2;
	(*processes)[2] = 1;
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

MOCK_FUNCTION_BEGIN(lmctfy_container_unregister_notification,
                    struct container *container,
                    const notification_id_t notif_id) {
} MOCK_FUNCTION_END
