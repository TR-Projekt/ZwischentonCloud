<h1 align="center">
    ZwischentonCloud
</h1>

<p align="center">
   <a href="https://github.com/TR-Projekt/zwischentoncloud/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/TR-Projekt/zwischentoncloud?style=flat"></a>
   <a href="https://github.com/TR-Projekt/zwischentoncloud/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/TR-Projekt/zwischentoncloud?style=flat"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/TR-Projekt/zwischentoncloud.svg"></a>
</p>

<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#usage">Usage</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#engage">Engage</a> •
  <a href="#licensing">Licensing</a>
</p>

TBA

## Development

TBA

### Requirements

TBA

### Setup

TBA

## Deployment

### VM deployment

The install, update and uninstall scripts should work with any system that uses *systemd* and *firewalld*.
Additionally the scripts will somewhat work under macOS but won't configure the firewall or launch service.

Installing
```bash
curl -o install.sh https://raw.githubusercontent.com/TR-Projekt/zwischentoncloud/main/operation/install.sh
chmod +x install.sh
sudo ./install.sh
```
Updating
```bash
curl -o update.sh https://raw.githubusercontent.com/TR-Projekt/zwischentoncloud/main/operation/update.sh
chmod +x update.sh
sudo ./update.sh
```

### Server

All of the scripts require Ubuntu 20 LTS as the operating system and that the server has already been initialised, see the steps to do that [here](https://github.com/Festivals-App/festivals-documentation/tree/master/deployment/general-vm-setup).

TBA

### Docker

TBA

## Usage

TBA

base/health
base/version
base/info
base/log

discovery.base/services
discovery.base/loversear

api.base/*

files.base/*


### Documentation

The zwischentoncloud is documented in detail [here](./DOCUMENTATION.md).

The full documentation for the Festivals App is in the [zwischenton-documentation](https://github.com/TR-Projekt/zwischenton-documentation) repository. 
The documentation repository contains technical documents, architecture information, UI/UX specifications, and whitepapers related to this implementation.

## Engage

TBA

The following channels are available for discussions, feedback, and support requests:

| Type                     | Channel                                                |
| ------------------------ | ------------------------------------------------------ |
| **General Discussion**   | <a href="https://github.com/TR-Projekt/zwischenton-documentation/issues/new/choose" title="General Discussion"><img src="https://img.shields.io/github/issues/TR-Projekt/zwischenton-documentation/question.svg?style=flat-square"></a> </a>   |
| **Concept Feedback**    | <a href="https://github.com/TR-Projekt/zwischenton-documentation/issues/new/choose" title="Open Concept Feedback"><img src="https://img.shields.io/github/issues/TR-Projekt/zwischenton-documentation/architecture.svg?style=flat-square"></a>  |
| **Other Requests**    | <a href="mailto:phisto05@gmail.com" title="Email Zwischenton Team"><img src="https://img.shields.io/badge/email-Zwischenton%20team-green?logo=mail.ru&style=flat-square&logoColor=white"></a>   |

## Licensing

Copyright (c) 2024 Simon Gaus.

Licensed under the **GNU Lesser General Public License v3.0** (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at https://www.gnu.org/licenses/lgpl-3.0.html.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the [LICENSE](./LICENSE) for the specific language governing permissions and limitations under the License.