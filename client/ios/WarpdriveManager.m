//
//  WarpdriveManager.m
//
//  Created by Ali Najafizadeh.
//  Copyright © 2017 Pressly Inc. All rights reserved.
//

#import <React/RCTBridgeModule.h>
#import <React/RCTRootView.h>

#import "WarpdriveManager.h"
#import "Warpdrive.framework/Versions/A/Headers/Warpdrive.h"

typedef RCTBridge*(^BridgeCallback)(void);
static BridgeCallback getBridge = nil;

@implementation WarpdriveManager

RCT_EXPORT_MODULE();

// bridge will be initialized by internal react-native
@synthesize bridge = _bridge;

- (instancetype) init {
    self = [super init];
    if (self) {
        getBridge = ^RCTBridge*{
            return self.bridge;
        };
    }
    return self;
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

+ (NSURL*)sourceBundleForApp:(NSString*)app
                andRolloutAt:(NSString*)rolloutAt
                andGroupName:(NSString*)groupName
               andServerAddr:(NSString*)serverAddr
               andDeviceCert:(NSString*)deviceCert
                andDeviceKey:(NSString*)deviceKey
                       andCA:(NSString*)caCert {
    
    NSBundle* bundle = [NSBundle mainBundle];
    
    NSString* bundleVersion = [bundle objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    NSString* bundlePath = [[NSBundle mainBundle] bundlePath];
    NSString* documentPath = [WarpdriveManager documentPathWithGroupName:groupName];
    NSString* platform = @"ios";
    NSError* error = nil;
    
    deviceCert = [[bundle URLForResource:deviceCert withExtension:@"crt"] absoluteString];
    deviceKey = [[bundle URLForResource:deviceKey withExtension:@"key"] absoluteString];
    caCert = [[bundle URLForResource:caCert withExtension:@"crt"] absoluteString];
    
    WarpdriveInit(bundlePath, documentPath, platform, app, rolloutAt, bundleVersion, serverAddr, deviceCert, deviceKey, caCert, &error);
    if (error != nil) {
        NSString* value = [error localizedDescription];
        NSLog(@"%@", value);
    }
    
    NSString* path = WarpdriveBundlePath();
    if ([path isEqualToString:@""]) {
        return nil;
    }
    
    return [NSURL URLWithString:[path stringByAddingPercentEncodingWithAllowedCharacters:[NSCharacterSet URLFragmentAllowedCharacterSet]]];
}

+ (void)reloadWithPath:(NSString*)path {
    RCTBridge* bridge = getBridge();
    [bridge setValue:[NSURL URLWithString:path] forKey:@"bundleURL"];
    [bridge reload];
}

RCT_EXPORT_METHOD(isAnyUpdate:(RCTResponseSenderBlock)callback) {
    NSString* release = WarpdriveIsAnyUpdate();
    if (release != nil) {
        callback(@[release]);
    } else {
        callback(@[[NSNull null]]);
    }
}

RCT_EXPORT_METHOD(update:(RCTResponseSenderBlock)callback) {
    NSError *err;
    WarpdriveUpdate(&err);
    if (err != nil) {
        callback(@[[err localizedDescription]]);
    } else {
        callback(@[[NSNull null]]);
    }
}

RCT_EXPORT_METHOD(reload) {
    NSString *path = WarpdriveBundlePath();
    if (path != nil && ![path isEqualToString:@""]) {
        [WarpdriveManager reloadWithPath:path];
    }
}


@end
