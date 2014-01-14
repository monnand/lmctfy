#ifndef SRC_CLI_COMMANDS_PAUSE_H_
#define SRC_CLI_COMMANDS_PAUSE_H_

#include <string>
using ::std::string;
#include <vector>

#include "util/task/status.h"

namespace containers {
namespace lmctfy {

class ContainerApi;

namespace cli {

class OutputMap;

::util::Status PauseContainer(const ::std::vector<string> &argv,
                              const ContainerApi *lmctfy,
                              ::std::vector<OutputMap> *output);

void RegisterPauseCommand();

} // namespace cli
} // namespace lmctfy
} // namespace containers

#endif // SRC_CLI_COMMANDS_PAUSE_H_
