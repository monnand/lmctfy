// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "lmctfy/cli/commands/pause.h"

#include <memory>
#include <string>
using ::std::string;
#include <vector>
#include <iostream>

#include "util/errors.h"
#include "strings/stringpiece.h"
#include "lmctfy/cli/command.h"
#include "util/task/statusor.h"

using ::std::unique_ptr;
using ::std::vector;
using ::util::Status;
using ::util::StatusOr;

namespace containers {
namespace lmctfy {
namespace cli {

// Command to resume a container.
Status ResumeContainer(const vector<string> &argv, const ContainerApi *lmctfy,
                       vector<OutputMap> *output) {
  // Args: pause <container name>
  if (argv.size() != 2) {
    return Status(::util::error::INVALID_ARGUMENT,
                  "See help for supported options.");
  }
  const string container_name = argv[1];
  unique_ptr<Container> container;
  RETURN_IF_ERROR(lmctfy->Get(container_name), &container);
  RETURN_IF_ERROR(container->Resume());
  return Status::OK;
}

void RegisterResumeCommand() {
  RegisterRootCommand(
      CMD("resume",
          "resume all processes in a container",
          "",
          CMD_TYPE_SETTER, 1, 1, &ResumeContainer));
}

} // namespace containers
} // namespace lmctfy
} // namespace cli

