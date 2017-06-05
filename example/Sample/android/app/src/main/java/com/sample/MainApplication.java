package com.sample;

import android.app.Application;
import android.support.annotation.Nullable;

import com.facebook.react.ReactApplication;
import com.facebook.react.ReactNativeHost;
import com.facebook.react.ReactPackage;
import com.facebook.react.shell.MainReactPackage;
import com.facebook.soloader.SoLoader;
import com.pressly.warpdrive.WarpdrivePackage;

import java.util.Arrays;
import java.util.List;

public class MainApplication extends Application implements ReactApplication {

  private final ReactNativeHost mReactNativeHost = new ReactNativeHost(this) {
    @Override
    public boolean getUseDeveloperSupport() {
      return BuildConfig.DEBUG;
    }

    @Override
    protected List<ReactPackage> getPackages() {
      String app = "sample";
      String rollout = "dev";
      String serverAddr = "192.168.1.183:10001";
      String deviceCert = "device.crt";
      String deviceKey = "device.key";
      String caCert = "ca.crt";

      return Arrays.<ReactPackage>asList(
          new MainReactPackage(),
          new WarpdrivePackage(MainApplication.this, app, rollout, serverAddr, deviceCert, deviceKey, caCert)
      );
    }

    @Override
    protected @Nullable String getJSBundleFile() {
      return WarpdrivePackage.sourceBundle();
    }
  };

  @Override
  public ReactNativeHost getReactNativeHost() {
    return mReactNativeHost;
  }

  @Override
  public void onCreate() {
    super.onCreate();
    SoLoader.init(this, /* native exopackage */ false);
  }
}
