//
//  EventCallback.m
//  Warpify
//
//  Created by Ali Najafizadeh on 2016-12-20.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//

#import "EventCallbackWrapper.h"

@implementation EventCallbackWrapper{
  EventCallbackBlock _block;
}

- (void)setBlock: (EventCallbackBlock)block{
  self->_block = block;
}

- (void)do:(long)kind value:(NSString*)value; {
  self->_block(kind, value);
}

@end
