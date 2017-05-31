//
//  WarpdriveManager.m
//
//  Created by Ali Najafizadeh.
//  Copyright © 2017 Pressly Inc. All rights reserved.
//

#import "WarpdriveManager.h"
#import "Warpdrive.framework/Headers/Warpdrive.h"

@implementation WarpdriveManager

RCT_EXPORT_MODULE();

// bridge will be initialized by internal react-native
@synthesize bridge = _bridge;

+ (NSURL*)sourceBundleForApp:(NSString*)app
                andRolloutAt:(NSString*)rolloutAt
               andServerAddr:(NSString*)serverAddr
               andDeviceCert:(NSString*)deviceCert
                andDeviceKey:(NSString*)deviceKey
                       andCA:(NSString*)cert {
    
}

@end
