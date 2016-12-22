//
//  WarpifyManager.m
//  Warpdrive
//
//  Created by Ali Najafizadeh on 2016-12-05.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "WarpifyManager.h"

static Warpify* _warpify = nil;

@implementation WarpifyManager

RCT_EXPORT_MODULE();

@synthesize bridge = _bridge;

RCT_EXPORT_METHOD(reload:(NSString*)path) {
  dispatch_async(dispatch_get_main_queue(), ^{
    [_warpify reloadFromPath:path];
  });
}

+ (NSURL*)sourceBundleWithDefaultCycle:(NSString*)defaultCycle groupName:(NSString*)groupName forceUpdate:(BOOL)forceUpdate {
  _warpify = [Warpify createWithDefaultCycle:defaultCycle forceUpdate:forceUpdate groupName:groupName];
  return [_warpify sourceBundle];
}

- (instancetype) init {
  self = [super init];
  if (self) {
    // we need to pass the _bridge to
    // warpify object so it can be called to reload the app
    [_warpify setBridgeCallback:^RCTBridge*{
      return self.bridge;
    }];
  }
  return self;
}

@end
