//
//  Warpdrive.m
//  Warpdrive
//
//  Created by Ali Najafizadeh on 2016-12-05.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "Warpify.h"
#import "Warpify.framework/Headers/Warpify.h"

#import "RCTBundleURLProvider.h"

@implementation Warpify

+ (instancetype)shared {
  static Warpify *sharedInstance;
  static dispatch_once_t once_token;
  
  dispatch_once(&once_token, ^{
    sharedInstance = [Warpify new];
  });
  
  return sharedInstance;
}

- (NSURL *)sourceBundle {
  return [[RCTBundleURLProvider sharedSettings] jsBundleURLForBundleRoot:@"index.ios" fallbackResource:nil];
}

- (void) call {
  
}

@end
