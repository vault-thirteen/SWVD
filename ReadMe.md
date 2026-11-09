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

`list.exe <AuthorName> <PagesCount> <CookiesFilePath> <DataFolder>`

`download.exe <AuthorName> <CookiesFilePath> <DataFolder> <Quality>`

### Example

> list.exe John123 2 settings\cookies.json data`

> download.exe John123 settings\cookies.json data 1080`
