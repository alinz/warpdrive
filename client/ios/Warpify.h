//
//  Warpify.h
//  Warpdrive
//
//  Created by Ali Najafizadeh on 2016-12-05.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import <Foundation/Foundation.h>

#import "RCTBridgeModule.h"

@interface Warpify : NSObject

// if you plan to update both share extension and app itself, you have to pass the groupName
// this is importnt because in ios, group directory is the one that can be access by both app and share extension
+ (instancetype)createWithDefaultCycle:(NSString*)defaultCycle forceUpdate:(BOOL)forceUpdate groupName:(NSString*)groupName;
// if you don't have any share extension and you want to update, use this method.
+ (instancetype)createWithDefaultCycle:(NSString*)defaultCycle forceUpdate:(BOOL)forceUpdate;
// shared will call createWithDefaultCycle with default value
// if you want to change the default value, you have to call createWithDefaultCycle first.
// calling shared will return the same object
+ (instancetype)shared;

- (NSURL *)sourceBundle;
- (void)reloadFromPath:(NSString*)path;

@end

