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

## Usage
| :exclamation:  Once you execute the below commands make sure to leave the Google Chrome browser that it launches open   |
|-----------------------------------------|

### Windows
| :memo:        | All commands should be executed inside of cmd.exe |
|---------------|:------------------------|
```Batchfile
./nvidia-clerk-windows.exe -region=REGION_CODE_HERE
```

### Mac OSX
| :memo:        | All commands should be executed inside of Terminal.app |
|---------------|:------------------------|
```shell
chmod +x ./nvidia-clerk-darwin

./nvidia-clerk-darwin -region=REGION_CODE_HERE
```

### Linux
| :memo:        | All commands should be executed inside of Shell |
|---------------|:------------------------|
```shell
chmod +x ./nvidia-clerk-linux

./nvidia-clerk-linux -region=REGION_CODE_HERE
```

## Testing

Testing is currenly only supported for the USA region but it should show you what the automated checkout will look like.

### Windows
| :memo:        | All commands should be executed inside of cmd.exe |
|---------------|:------------------------|
```Batchfile
./nvidia-clerk-windows.exe -region=USA -test
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
./nvidia-clerk-windows.exe -sms -test
```

### Usage

```shell
./nvidia-clerk-windows.exe -sms -region=REGION_CODE_HERE
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
./nvidia-clerk-windows.exe -discord -region=REGION_CODE_HERE
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
./nvidia-clerk-windows.exe -twitter -region=REGION_CODE_HERE
```
