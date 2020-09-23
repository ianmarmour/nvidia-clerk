# nvidia-clerk

This project was written in response to the recent NVIDIA RTX 3080 release debacle. During the launch multiple different groups of scalpers used
closed source "bots" to procure large quantities of NVIDIA GPU's and most consumers were left without being able to purchase the product. This 
project will provide a short term solution so that customers can ensure they can buy a GPU and compete with these scalpers.

NVIDIA Clerk doesn't actually purchase products for customers, it simply tracks the avaliable inventory from NVIDIAs APIs and automatically adds a GPU
to your checkout/cart and navigates your browser checkout page whenever they become avaliable. The clerk can also notify you of this process if you provide
Twilio API information (I am not interested in running an entire service for users, so this feature is limited to users aware of how to setup
such things).

## Install

### Requirements

* Google Chrome
* Administrator Access

### Download
Download the latest release from [Releases](https://github.com/ianmarmour/nvidia-clerk/releases/latest)
| :exclamation:  Make sure you accept any browser warnings, these warnings are due to the fact that these release binaries are not "signed" (this costs money and as a free project we haven't paid for a signing certificate)   |
|-----------------------------------------|

### Supported Region Codes

AUT,BEL,CAN,CZE,DNK,FIN,FRA,DEU,USA,GBR,IRL,ITA,SWE,LUX,POL,PRT,ESP, NOR, NLD

### Supported Models

3070, 3080, 3090 

## Usage
| :exclamation:  Once you execute the below commands make sure to leave the Google Chrome browser that it launches open   |
|-----------------------------------------|

### Windows
| :memo:        | All commands should be executed inside of cmd.exe |
|---------------|:------------------------|
```Batchfile
nvidia-clerk-windows.exe -region=REGION_CODE_HERE -model=3080
```

### Mac OSX
| :memo:        | All commands should be executed inside of Terminal.app |
|---------------|:------------------------|
```shell
chmod +x ./nvidia-clerk-darwin

./nvidia-clerk-darwin -region=REGION_CODE_HERE -model=3080
```

### Linux
| :memo:        | All commands should be executed inside of Shell |
|---------------|:------------------------|
```shell
chmod +x ./nvidia-clerk-linux

./nvidia-clerk-linux -region=REGION_CODE_HERE -model=3080
```

## Testing

Testing is currenly only supported for the USA region but it should show you what the automated checkout will look like.

### Windows
| :memo:        | All commands should be executed inside of cmd.exe |
|---------------|:------------------------|
```Batchfile
nvidia-clerk-windows.exe -region=USA -test
```

### Mac OSX
| :memo:        | All commands should be executed inside of Terminal.app |
|---------------|:------------------------|
```shell
./nvidia-clerk-darwin -region=USA -test
```

### Linux
| :memo:        | All commands should be executed inside of Shell |
|---------------|:------------------------|
```shell
./nvidia-clerk-linux -region=USA -test
```


# Advanced Usage

## Manual Delay Usage
Example of setting a 1 second delay (delay is specificed in miliseconds)
```Batch
nvidia-clerk-windows.exe -region=USA -delay=1000
```

## SMS Notifications

### Configuration
```Batchfile
set TWILIO_ACCOUNT_SID=YOUR_TWILIO_ACCOUNT_SID_HERE
set TWILIO_TOKEN=YOUR_TWILIO_TOKEN_HERE
set TWILIO_SOURCE_NUMBER=YOUR_TWILIO_SERVICE_NUMBER_HERE
set TWILIO_DESTINATION_NUMBER=YOUR_DESITNATION_NUMBER_FOR_NOTIFICATIONS_HERE
```

### Testing
```shell
nvidia-clerk-windows.exe -sms -test
```

### Usage

```shell
nvidia-clerk-windows.exe -sms -region=REGION_CODE_HERE -model=3080
```

## Discord Notifications

### Configuration
```Batchfile
set DISCORD_WEBHOOK_URL=DISCORD_WEBHOOK_URL_HERE
```

### Testing
```Batchfile
./nvidia-clerk-windows.exe -discord -test
```

### Usage

```Batchfile
./nvidia-clerk-windows.exe -discord -region=REGION_CODE_HERE -model=3080
```

## Twitter Notifications

### Configuration
```Batchfile
set TWITTER_CONSUMER_KEY=YOUR_TWITTER_CONSUMER_KEY_HERE
set TWITTER_CONSUMER_SECRET=YOUR_TWITTER_CONSUMER_SECRET_HERE
set TWITTER_ACCESS_TOKEN=YOUR_TWITTER_ACCESS_TOKEN_HERE
set TWITTER_ACCESS_SECRET=YOUR_TWITTER_ACCESS_SECRET_HERE
```

### Testing
```Batchfile
./nvidia-clerk-windows.exe -twitter -test
```

### Usage

```Batchfile
./nvidia-clerk-windows.exe -twitter -region=REGION_CODE_HERE -model=3080
```

## Telegram Notifications

### Configuration
1. In Telegram, create a new bot by contacting the `@BotFather` and either creating a new bot or getting an existing bot token.
2. Once you have your token (e.g. `123456789:LaeQH7pvhJbqmqV9_akkVbClpsSZr7su63o`), add your bot to a group or channel.
3. Get the chat id of your group or channel by following the instructions [available on GitHub here](https://github.com/GabrielRF/telegram-id). This will be either in the format `-10012345678` or `@my_channel_or_group`.
4. Use the full token, including the `:`, as your API key, and use the full chat id, including the `-` or `@`, as your chat id.

```Batchfile
set TELEGRAM_API_KEY=YOUR_TELEGRAM_API_KEY_HERE
set TELEGRAM_CHAT_ID=YOUR_TELEGRAM_CHAT_ID_HERE
```

### Testing
```Batchfile
./nvidia-clerk-windows.exe -telegram -test
```

### Usage

```Batchfile
./nvidia-clerk-windows.exe -telegram -region=REGION_CODE_HERE -model=3080
```

## FAQ

#### exec: "google-chrome": executable file not found in %PATH%
The path to google chrome needs to be added to the path. While reinstalling Chrome can fix it, a quick search on Google will show you how to do exactly this.

#### My environment variables are being unset
You can add them to your .bashrc or .zsh files and be sure to export them. Sourcing the files and/or restarting your window will load them. That said, the easiest way to do this is just to wrap these values in a wrapper script. eg.

```
set TWILIO_ACCOUNT_SID={YOUR_TWILIO_ACCOUNT_SID}
set TWILIO_TOKEN={YOUR_TWILIO_TOKEN}
set TWILIO_SOURCE_NUMBER={YOUR_TWILIO_PHONE_NR}
set TWILIO_DESTINATION_NUMBER={THE_NUMBER_TO_WHICH_TO_SEND_THE_NOTIFICATIONS}

nvidia-clerk-windows.exe -region={REGION_CODE_HERE} -model=3080
``` 

#### The log says "Access Denied" and then quits
This can be caused by a myriad of things. However, here's some items to check:
1. Is your antivirus blocking it? Some AVs will copy the threatening file into a separate directory and give a notification
2. Is your firewall blocking access?
3. Do you have the right permissions to run this program?

#### It's still not working... halp
Please search in the Issues tab of this repository for solutions, and if you can't find anything, post a detailed listing. Please follow the following structure for asking questions:

```
Operating system: 
Version of nvidia-clerk you're running:
Parameters passed in (redact sensitive info):
Additional info:
```
