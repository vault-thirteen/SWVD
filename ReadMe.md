# SWVD

Simple Website Video Downloader.

The bundle consists of two executables:
- `list.exe` which lists author's videos;
- `download.exe` which downloads author's videos.

When requested video quality is not available, the nearest lower quality is 
used.

This downloader was created to be used for a specific website. It means that 
this tool does not work on all websites in the Internet network.  

Ideas from this tool can be taken to create more complex tools. 

## Usage

`list.exe <SettingsFilePath> <AuthorName> <PagesCount>`

`download.exe <SettingsFilePath> <AuthorName> <Quality>`

### Example

> list.exe settings\settings.json Author 2

> download.exe settings\settings.json Author 1080

## Installation

`go install github.com/vault-thirteen/SWVD/cmd/list@latest`

`go install github.com/vault-thirteen/SWVD/cmd/download@latest`
