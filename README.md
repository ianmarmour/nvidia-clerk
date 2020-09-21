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

## Determine Your SKU

In order to configure the nvidia-clerk you'll need to determine which SKU you would like the clerk to monitor (currently the clerk only supports monitoring a
single SKU per an instance)

### Supported Region SKUs

**United States**
| Product Name | SKU |
|---|---|
| Nvidia RTX 3070 FE  | N/A |
| Nvidia RTX 3080 FE  | 5438481700 |
| Nvidia RTX 3090 FE  | N/A |

**Great Britain**
| Product Name | SKU |
|---|---|
| Nvidia RTX 3070 FE  | N/A |
| Nvidia RTX 3080 FE  | 5438792800 |
| Nvidia RTX 3090 FE  | N/A |

**Germany**
| Product Name | SKU |
|---|---|
| Nvidia RTX 3070 FE  | N/A |
| Nvidia RTX 3080 FE  | 5438792300 |
| Nvidia RTX 3090 FE  | N/A |

**France**
| Product Name | SKU |
|---|---|
| Nvidia RTX 3070 FE  | N/A |
| Nvidia RTX 3080 FE  | 5438795200 |
| Nvidia RTX 3090 FE  | N/A |

**Sweden**
| Product Name | SKU |
|---|---|
| Nvidia RTX 3070 FE  | N/A |
| Nvidia RTX 3080 FE  | 5438798100 |
| Nvidia RTX 3090 FE  | N/A |


## Configuration

### Windows
| :memo:        | All commands should be executed inside of cmd.exe |
|---------------|:------------------------|
```Batchfile
set NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
```

### Mac OSX
| :memo:        | All commands should be executed inside of Terminal.app |
|---------------|:------------------------|
```shell
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
```

### Linux
| :memo:        | All commands should be executed inside of Shell |
|---------------|:------------------------|
```shell
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
```

## Usage
| :exclamation:  Once you execute the below commands make sure to leave the Google Chrome browser that it launches open   |
|-----------------------------------------|

### Windows
| :memo:        | All commands should be executed inside of cmd.exe |
|---------------|:------------------------|
```Batchfile
./nvidia-clerk-windows.exe
```

### Mac OSX
| :memo:        | All commands should be executed inside of Terminal.app |
|---------------|:------------------------|
```shell
chmod +x ./nvidia-clerk-darwin

./nvidia-clerk-darwin
```

### Linux
| :memo:        | All commands should be executed inside of Shell |
|---------------|:------------------------|
```shell
chmod +x ./nvidia-clerk-linux

./nvidia-clerk-linux
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
./nvidia-clerk-windows.exe -sms
```

## Unsupported Regions

You can attempt to manually configure you're sku, locale and currency options using the following steps support for your region may be limited, if you manage to get a new region working please cut a GitHub issue and include the SKU, LOCALE and Currency settings that you used so we can add first party support for this region.

### Determining Your Locale

US readers can ignore this section however if you are non-us please determine your locale from [this list](https://www.science.co.il/language/Locale-codes.php), make sure you format your local as xx_xx vs xx-xx otherwise the bot will fail.

### Determining Your Currency Code

US readers can ignore this section however if you are non-us please determine your currency code from [this list](https://www.iban.com/currency-codes), make sure you use the 3 letter currency code for your region from this list otherwise the bot will fail.

### Determining Your SKU 

For more advanced users attempting to identify the SKU for their countris store please update the below URL's locale to your respective locale and iterate through the pageNumbers by updating the page number at the end of the URL until you find the 3070/3080/3090 SKU number for your region (NOTE: NOT ALL REGIONS HAVE THE SKUS AVALIABLE YET)

https://api.digitalriver.com/v1/shoppers/me/products?apiKey=9485fa7b159e42edb08a83bde0d83dia&locale=en_ca&format=json&expand=product&fields=product.id,product.displayName,product.pricing&pageNumber=1

### Additional Configuration

After you have your new SKU, Currency Code and Locale please set the configuration up as follows,

#### Windows

```Batchfile
set NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
set NVIDIA_CLERK_LOCALE=YOUR_LOCALE_HERE
set NVIDIA_CLERK_CURRENCY=YOUR_CURRENCY_HERE
```

#### Mac

```shell
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
export NVIDIA_CLERK_LOCALE=YOUR_LOCALE_HERE
export NVIDIA_CLERK_CURRENCY=YOUR_CURRENCY_HERE
```

#### Linux

```shell
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
export NVIDIA_CLERK_LOCALE=YOUR_LOCALE_HERE
export NVIDIA_CLERK_CURRENCY=YOUR_CURRENCY_HERE
```

# Troubleshooting

## Windows

#### Chrome install path configuration

If you encounter an error about your chrome path but it is installed, set the `CHROME_PATH` env var to the location of its exe.
For example "C:/Program Files/Google/Chrome/Application/chrome.exe".
You can find the path by right clicking on chrome in your start menu and choosing "Open File Location", right click the
.exe file, select 'Properties', and copy and paste the Target.

```Batchfile
set NVIDIA_CLERK_CHROME_PATH=PATH/TO/YOUR/CHROME.EXE
```
