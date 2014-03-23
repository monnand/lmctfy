// Mock functions for clmctfy
#include "clmctfy-raw.h"
#include <stdlib.h>
#include "clmctfy.h"
#include "clmctfy_mock.h"
#include "macros.h"

static notification_id_t next_notif_id = 0;

MOCK_FUNCTION_BEGIN(lmctfy_init_machine_raw, const void *spec, const size_t spec_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_new_container_api, struct container_api **api) {
	// Some random pointer.
	// Will not be dereferenced.
	*api = CONTAINER_API_ADDR;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_api_create_container_raw,
    struct container_api *api,
    const char *container_name,
    const void *spec,
    const size_t spec_size,
    struct container **container) {
	*container = (struct container *)malloc(sizeof(struct container));
	(*container)->name = strdup(container_name);
	(*container)->callback = NULL;
	(*container)->cb_userdata = NULL;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_run_raw,
                    struct container *container,
                    const int argc,
                    const char **argv,
                    const void *spec,
                    const size_t spec_size,
                    pid_t *tid) {
	*tid = 10;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_update_raw,
                    struct container *container,
                    int policy,
                    const void *spec,
                    const size_t spec_size) {
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_register_notification_raw,
                    struct container *container,
                    lmctfy_event_callback_f callback,
                    void *user_data,
                    const void *spec,
                    const size_t spec_size,
                    notification_id_t *notif_id) {
	*notif_id = next_notif_id++;
	container->callback = callback;
	container->cb_userdata = user_data;
} MOCK_FUNCTION_END

MOCK_FUNCTION_BEGIN(lmctfy_container_stats_raw,
                    struct container *container,
                    int stats_type,
                    void **stats,
                    size_t *stats_size) {
	container->callback = NULL;
	container->cb_userdata = NULL;
} MOCK_FUNCTION_END
