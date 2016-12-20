//
//  ReloadTaskWrapper.m
//  Warpify
//
//  Created by Ali Najafizadeh on 2016-12-20.
//  Copyright © 2016 Ali Najafizadeh. All rights reserved.
//
#import "ReloadTaskWrapper.h"

@implementation ReloadTaskWrapper{
  ReloadBlock _block;
}

- (void)setBlock: (ReloadBlock)block{
  self->_block = block;
}

- (void)execute:(NSString*)path {
  self->_block(path);
}

@end
