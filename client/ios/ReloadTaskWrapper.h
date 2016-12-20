//
//  ReloadTaskWrapper.h
//  Sample1
//
//  Created by Ali Najafizadeh on 2016-12-20.
//  Copyright © 2016 Facebook. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "Warpify.framework/Headers/Warpify.h"

typedef void (^ReloadBlock)(NSString* path);

@interface ReloadTaskWrapper : NSObject<GoWarpifyTask>
- (void)setBlock: (ReloadBlock)block;
@end
