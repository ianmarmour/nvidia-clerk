# nvidia-clerk

This project was written in response to the recent NVIDIA RTX 3080 release debacle. During the launch multiple different groups of scalpers used
closed source "bots" to procure large quantities of NVIDIA GPU's and most consumers were left without being able to purchase the product. This 
project attempts will provide a short term solution so that customers can ensure they can manually buy a GPU and actually compete with these scalpers.

NVIDIA Clerk doesn't actually purchase products for customers, it simply tracks the avaliable inventory from NVIDIAs APIs and automatically adds a card
to your cart and drops your browser to the checkout page whenever they become avaliable. The clerk can also notify you of this process if you provide
Twilio API information (I currently am not interested in running an entire service for users, so this feature is limited to users aware of how to setup
such things).

## Usage

### SKUs

In order to configure the nvidia-clerk you'll need to determine which SKU you would like the clerk to monitor (currently the clerk only supports monitoring a
single SKU per an instance)

| Product Name | SKU |
|---|---|
| Nvidia RTX 3070 FE  | N/A |
| Nvidia RTX 3080 FE  | 5438481700 |
| Nvidia RTX 3090 FE  | N/A |


### Windows

#### Without SMS Notification

##### Setup Configuration
```
set NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
```

##### Run Nvidia Clerk
```
nvidia-clerk.exe
```

#### With SMS Notification

##### Setup Configuration
```
set NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
set TWILIO_ACCOUNT_SID=YOUR_TWILIO_ACCOUNT_SID_HERE
set TWILIO_SERVICE_SID=YOUR_TWILIO_SERVICE_SID_HERE
set TWILIO_TOKEN=YOUR_TWILIO_TOKEN_HERE
set TWILIO_SOURCE_NUMBER=YOUR_TWILIO_SERVICE_NUMBER_HERE
set TWILIO_DESTINATION_NUMBER=YOUR_DESITNATION_NUMBER_FOR_NOTIFICATIONS_HERE
```

##### Run Nvidia Clerk
```
nvidia-clerk.exe -sms
```

### Mac OSX

#### Without SMS Notification

##### Setup Configuration
```
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
```

##### Run Nvidia Clerk
```
chmod +x ./nvidia-clerk

./nvidia-clerk
```

#### With SMS Notification

##### Setup Configuration
```
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
export TWILIO_ACCOUNT_SID=YOUR_TWILIO_ACCOUNT_SID_HERE
export TWILIO_SERVICE_SID=YOUR_TWILIO_SERVICE_SID_HERE
export TWILIO_TOKEN=YOUR_TWILIO_TOKEN_HERE
export TWILIO_SOURCE_NUMBER=YOUR_TWILIO_SERVICE_NUMBER_HERE
export TWILIO_DESTINATION_NUMBER=YOUR_DESITNATION_NUMBER_FOR_NOTIFICATIONS_HERE
```

##### Run Nvidia Clerk
```
chmod +x ./nvidia-clerk

./nvidia-clerk -sms
```

### Linux

#### Without SMS Notification

##### Setup Configuration
```
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
```

##### Run Nvidia Clerk
```
chmod +x ./nvidia-clerk

./nvidia-clerk
```

#### With SMS Notification

##### Setup Configuration
```
export NVIDIA_CLERK_SKU=YOUR_DESIRED_SKU_HERE
export TWILIO_ACCOUNT_SID=YOUR_TWILIO_ACCOUNT_SID_HERE
export TWILIO_SERVICE_SID=YOUR_TWILIO_SERVICE_SID_HERE
export TWILIO_TOKEN=YOUR_TWILIO_TOKEN_HERE
export TWILIO_SOURCE_NUMBER=YOUR_TWILIO_SERVICE_NUMBER_HERE
export TWILIO_DESTINATION_NUMBER=YOUR_DESITNATION_NUMBER_FOR_NOTIFICATIONS_HERE
```

##### Run Nvidia Clerk
```
chmod +x ./nvidia-clerk

./nvidia-clerk -sms
```
