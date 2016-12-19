//
//  Warpify.m
//  Warpdrive
//
//  Created by Ali Najafizadeh on 2016-12-05.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "Warpify.h"
#import "Warpify.framework/Headers/Warpify.h"

#import "RCTBundleURLProvider.h"

static Warpify *sharedInstance;

@implementation Warpify

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
  
  dispatch_once(&once_token, ^{
    sharedInstance = [Warpify new];
    
    // need to update the internal go variable
    NSString* bundleVersion = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    NSString* bundlePath = [[NSBundle mainBundle] bundlePath];
    NSString* documentPath = [Warpify documentPathWithGroupName:groupName];
    NSString* platform = @"ios";
    
    GoWarpifySetup(bundleVersion, bundlePath, documentPath, platform, defaultCycle, forceUpdate);
  });
  
  return sharedInstance;
}

+ (instancetype)createWithDefaultCycle:(NSString*)defaultCycle forceUpdate:(BOOL)forceUpdate {
  return [Warpify createWithDefaultCycle:defaultCycle forceUpdate:forceUpdate groupName:nil];
}

+ (instancetype)shared {
  return [Warpify createWithDefaultCycle:@"prod" forceUpdate:false];
}

- (NSURL *)sourceBundle {
  return [[RCTBundleURLProvider sharedSettings] jsBundleURLForBundleRoot:@"index.ios" fallbackResource:nil];
}

@end
