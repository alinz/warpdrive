//
//  Warpdrive.m
//  Warpdrive
//
//  Created by Ali Najafizadeh on 2016-12-05.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "Warpdrive.h"
#import "Warpdrive.framework/Headers/Warpdrive.h"

#import "RCTBundleURLProvider.h"

@implementation Warpdrive

+ (instancetype)shared {
  static Warpdrive *sharedInstance;
  static dispatch_once_t once_token;
  
  dispatch_once(&once_token, ^{
    sharedInstance = [Warpdrive new];
  });
  
  return sharedInstance;
}

- (NSURL *)sourceBundle {
  return [[RCTBundleURLProvider sharedSettings] jsBundleURLForBundleRoot:@"index.ios" fallbackResource:nil];
}

- (void) call {
  GoWarpdriveSetup(@"Hello", @"World");
}

@end
