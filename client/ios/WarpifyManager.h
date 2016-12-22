//
//  WarpifyManager.h
//  Warpify
//
//  Created by Ali Najafizadeh on 2016-12-21.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "Warpify.h"

@interface WarpifyManager : NSObject <RCTBridgeModule>
+ (NSURL*)sourceBundleWithDefaultCycle:(NSString*)defaultCycle groupName:(NSString*)groupName forceUpdate:(BOOL)forceUpdate;
@end
