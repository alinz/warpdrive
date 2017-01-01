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

// bridge will be initialized by internal react-native
@synthesize bridge = _bridge;

+ (NSURL*)sourceBundleWithDefaultCycle:(NSString*)defaultCycle groupName:(NSString*)groupName forceUpdate:(BOOL)forceUpdate {
  _warpify = [Warpify createWithDefaultCycle:defaultCycle forceUpdate:forceUpdate groupName:groupName];
  return [_warpify sourceBundle];
}

- (instancetype) init {
  self = [super init];
  if (self) {
    // we need to pass the _bridge to
    // warpify object so it can be called to reload the app
    // I can't find a way to make sure the bridge is initialized before
    // being used.
    [_warpify setBridgeCallback:^RCTBridge*{
      return self.bridge;
    }];
  }
  return self;
}

////////////////////////////////////////
// exported methods to javascript
////////////////////////////////////////

RCT_EXPORT_METHOD(cycles:(RCTResponseErrorBlock)errCallback
         successCallback:(RCTResponseSenderBlock)successCallback) {
  NSError* error;
  NSString* result = [_warpify cyclesWithError:&error];
  if (error != nil) {
    errCallback(error);
  } else {
    successCallback(@[result]);
  }
}

RCT_EXPORT_METHOD(remoteVersions:(int)cycleId
                     errCallback:(RCTResponseErrorBlock)errCallback
                 successCallback:(RCTResponseSenderBlock)successCallback) {
  NSError* error;
  NSString* result = [_warpify remoteVersionsWithCycleId:cycleId error:&error];
  if (error != nil) {
    errCallback(error);
  } else {
    successCallback(@[result]);
  }
}

RCT_EXPORT_METHOD(localVersions:(int)cycleId
                    errCallback:(RCTResponseErrorBlock)errCallback
                successCallback:(RCTResponseSenderBlock)successCallback) {
  NSError* error;
  NSString* result = [_warpify localVersionsWithCycleId:cycleId error:&error];
  if (error != nil) {
    errCallback(error);
  } else {
    successCallback(@[result]);
  }
}

RCT_EXPORT_METHOD(latestVersion:(int)cycleId
                    errCallback:(RCTResponseErrorBlock)errCallback
                successCallback:(RCTResponseSenderBlock)successCallback) {
  NSError* error;
  NSString* result = [_warpify latestVersionWithCycleId:cycleId error:&error];
  if (error != nil) {
    errCallback(error);
  } else {
    successCallback(@[result]);
  }
}

RCT_EXPORT_METHOD(downloadVersion:(int)cycleId
                          version:(NSString*)version
                      errCallback:(RCTResponseErrorBlock)errCallback
                  successCallback:(RCTResponseSenderBlock)successCallback) {
  NSError* error;
  [_warpify downloadVersionWithCycleID:cycleId version:version error:&error];
  if (error != nil) {
    errCallback(error);
  } else {
    successCallback(@[]);
  }
}

RCT_EXPORT_METHOD(reload:(int)cycleId
                 version:(NSString*)version
             errCallback:(RCTResponseErrorBlock)errCallback
         successCallback:(RCTResponseSenderBlock)successCallback) {
  NSError* error;
  [_warpify reloadVersionWithCycleID:cycleId version:version error:&error];
  if (error != nil) {
    errCallback(error);
  } else {
    successCallback(@[]);
  }
}

// becuase all the warpify internal calls require to access network or
// disk, it is better to use a different queue
- (dispatch_queue_t)methodQueue {
  return dispatch_queue_create("com.pressly.warpify", DISPATCH_QUEUE_SERIAL);
}

//// download requests a version to download, it simply start download and store the bundles into the local document
//// folder. if anything goes wrong, it will call callback with error. callback(err), otherwise callback err will be null
//RCT_EXPORT_METHOD(download:(NSString*)version callback:(RCTResponseSenderBlock)callback) {
//  
//}
//
//// localVersions retuns the list of available versions on local
//RCT_EXPORT_METHOD(localVersions:(RCTResponseSenderBlock)callback) {
//  
//}
//
//// removeVersion returns the list of available versions based on newest first
//RCT_EXPORT_METHOD(remoteVersions:(RCTResponseSenderBlock)callback) {
//  
//}
//
//RCT_EXPORT_METHOD(reload:(NSString*)version) {
//  //  dispatch_async(dispatch_get_main_queue(), ^{
//  //    [_warpify reloadFromPath:path];
//  //  });
//}

@end
