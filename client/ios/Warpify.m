//
//  Warpify.m
//  Warpdrive
//
//  Created by Ali Najafizadeh on 2016-12-05.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "React/RCTBundleURLProvider.h"
#import "React/RCTBridgeModule.h"
#import "React/RCTConvert.h"
#import "React/RCTEventDispatcher.h"
#import "React/RCTRootView.h"
#import "React/RCTUtils.h"

#import "Warpify.framework/Headers/Warpify.h"

#import "Warpify.h"
#import "EventCallbackWrapper.h"

static Warpify *sharedInstance;

@implementation Warpify {
  // we are creating this variable to make sure it never
  // garbage collected
  EventCallbackWrapper* _reloadCallback;
  BridgeCallback _bridgeCallback;
}

// this method returns the document path based on whether groupName given or not
+ (NSString*)documentPathWithGroupName:(NSString*)groupName {
  if (groupName == nil) {
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    return [paths objectAtIndex:0];
  } else {
    NSURL* pathURL = [[NSFileManager defaultManager] containerURLForSecurityApplicationGroupIdentifier:groupName];
    return [pathURL path];
  }
}

+ (instancetype)createWithDefaultCycle:(NSString*)defaultCycle forceUpdate:(BOOL)forceUpdate groupName:(NSString*)groupName {
  static dispatch_once_t once_token;
  
  // we are going to call this method once so subsiqent call to `shared` and `createWithDefaultCycle`
  // returns the same instance
  dispatch_once(&once_token, ^{
    sharedInstance = [Warpify new];
    
    // need to update the internal go variable
    NSString* bundleVersion = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    NSString* bundlePath = [[NSBundle mainBundle] bundlePath];
    NSString* documentPath = [Warpify documentPathWithGroupName:groupName];
    NSString* platform = @"ios";

    // set the reload path here
    sharedInstance->_reloadCallback = [EventCallbackWrapper new];
    [sharedInstance->_reloadCallback setBlock:^(long kind, NSString* path) {
      
      // HACK, this is just a hack becuase, init in WarpifyManager won't be called sync
      // so we don't know when the bridge callback is going to set. by adding a little bit of delay
      // we are kind of waiting and hopefully the RCTBridgeModule calls the init method 
      dispatch_time_t delay = dispatch_time(DISPATCH_TIME_NOW, NSEC_PER_SEC * 2.0);
      dispatch_after(delay, dispatch_get_main_queue(), ^(void){
        [sharedInstance reloadFromPath:path];
      });
      
    }];
    WarpifySetReload((EventCallbackWrapper*)sharedInstance->_reloadCallback);
    
    // Setup the basic requirements
    NSError* err;
    WarpifySetup(bundleVersion, bundlePath, documentPath, platform, defaultCycle, forceUpdate, &err);
    
    // if some error happends, we are forcefully set sharedInstance to nil
    // to crash the app at the beginning.
    if (err) {
      NSLog(err);
      sharedInstance = nil;
    }
  });
  
  return sharedInstance;
}

+ (instancetype)createWithDefaultCycle:(NSString*)defaultCycle forceUpdate:(BOOL)forceUpdate {
  return [Warpify createWithDefaultCycle:defaultCycle forceUpdate:forceUpdate groupName:nil];
}

+ (instancetype)shared {
  return [Warpify createWithDefaultCycle:@"prod" forceUpdate:false];
}

- (void)setBridgeCallback:(BridgeCallback) bridgeCallback {
  _bridgeCallback = bridgeCallback;
}

- (NSURL *)sourceBundle {
  NSString* path = WarpifySourcePath();
  if (path == nil || [path isEqualToString:@""]) {
    return [[RCTBundleURLProvider sharedSettings] jsBundleURLForBundleRoot:@"index.ios" fallbackResource:nil];
  }
  
  return [NSURL URLWithString:path];
}

- (void)reloadFromPath:(NSString*)path {
  // Since the bridge will be set automatically, we don't know when,
  // so in order to get the ref to bridge, we need to call the block code
  // which returns the correct instance of bridge
  RCTBridge* bridge = _bridgeCallback();
  // bundleURL has to be NSURL, if you pass it as string it will blow out
  [bridge setValue:[NSURL URLWithString:path] forKey:@"bundleURL"];
  [bridge reload];
}

- (NSString*)cyclesWithError:(NSError**)err {
  return WarpifyCycles(err);
}

- (NSString*)remoteVersionsWithCycleId:(int64_t)cycleID error:(NSError**)err {
  return WarpifyRemoteVersions(cycleID, err);
}

- (NSString*)localVersionsWithCycleId:(int64_t)cycleID error:(NSError**)err {
  return WarpifyLocalVersions(cycleID, err);
}

- (NSString*)latestVersionWithCycleId:(int64_t)cycleID error:(NSError**)err {
  return WarpifyLatest(cycleID, err);
}

- (void)downloadVersionWithCycleID:(int64_t)cycleId version:(NSString*)version error:(NSError**)err {
  WarpifyDownloadVersion(cycleId, version, err);
}

- (void)reloadVersionWithCycleID:(int64_t)cycleId version:(NSString*)version error:(NSError**)err {
  WarpifyReload(cycleId, version, err);
}

@end
