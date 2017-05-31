//
//  WarpdriveManager.h
//
//  Created by Ali Najafizadeh.
//  Copyright © 2017 Pressly Inc. All rights reserved.
//

#import <Foundation/Foundation.h>

#import "React/RCTBridgeModule.h"

@interface WarpdriveManager : NSObject <RCTBridgeModule>
+ (NSURL*)sourceBundleForApp:(NSString*)app
                andRolloutAt:(NSString*)rolloutAt
               andServerAddr:(NSString*)serverAddr
               andDeviceCert:(NSString*)deviceCert
                andDeviceKey:(NSString*)deviceKey
                       andCA:(NSString*)cert;
@end
