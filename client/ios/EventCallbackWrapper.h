//
//  EventCallback.h
//  Warpify
//
//  Created by Ali Najafizadeh on 2016-12-20.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "Warpify.framework/Headers/Warpify.h"

typedef void (^EventCallbackBlock)(long, NSString* path);

@interface EventCallbackWrapper : NSObject<GoWarpifyCallback>
- (void)do:(long)kind value:(NSString*)value;
- (void)setBlock: (EventCallbackBlock)block;
@end
