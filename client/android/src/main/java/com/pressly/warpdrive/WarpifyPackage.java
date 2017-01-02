package com.pressly.warpdrive;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;

import com.facebook.react.ReactPackage;
import com.facebook.react.bridge.Callback;
import com.facebook.react.bridge.JavaScriptModule;
import com.facebook.react.bridge.NativeModule;
import com.facebook.react.bridge.ReactApplicationContext;
import com.facebook.react.bridge.ReactContextBaseJavaModule;
import com.facebook.react.bridge.ReactMethod;
import com.facebook.react.uimanager.ViewManager;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import go.warpify.Warpify;

public class WarpifyPackage implements ReactPackage {
    private static String sourceBundlePath;

    public static String sourceBundle() {
        String path = null;
        File file = new File(sourceBundlePath);
        if (file.exists()) {
            path = sourceBundlePath;
        }
        return path;
    }

    private static String getBundleVersion(Context context) {
        PackageManager manager = context.getPackageManager();
        PackageInfo info = null;
        try {
            info = manager.getPackageInfo(context.getPackageName(), 0);
            return info.versionName;
        } catch (PackageManager.NameNotFoundException e) {
            e.printStackTrace();
            return "";
        }
    }

    private static int copy(OutputStream des, InputStream src) throws IOException {
        byte[] buffer = new byte[1024]; // 1kb for now it's enough for us
        int total = 0;

        int read;
        while ((read = src.read(buffer)) != -1) {
            des.write(buffer, 0, read);
            total += read;
        }

        return total;
    }

    private static String joinPath(String a, String b) {
        return new File(a, b).getPath();
    }

    private static boolean copyWarpFile(Context context, String path) throws IOException {
        InputStream src = context.getAssets().open("WarpFile");

        File file = new File(path);
        file.mkdirs();

        path = joinPath(path, "WarpFile");

        OutputStream des = null;
        try {
            des = new FileOutputStream(path);
        } catch (IOException e) {
            src.close();
            throw e;
        }

        return copy(des, src) != 0;
    }

    public WarpifyPackage(Context context, String defaultCycle, boolean forceUpdate) {
        String documentPath = context.getFilesDir().getAbsolutePath();
        String bundlePath = joinPath(documentPath, "warptemp");
        String bundleVersion = getBundleVersion(context);
        String platform = "android";

        // since we can't get the bundle path in android,
        // we need to copy the WarpFile into document folder under `warptemp folder`
        // this file will be overwritten every time we launch the app
        try {
            if (!copyWarpFile(context, bundlePath)) {
                return;
            }
        } catch (IOException e) {
            e.printStackTrace();
            return;
        }

        try {
            Warpify.setup(bundleVersion, bundlePath, documentPath, platform, defaultCycle, forceUpdate);
        } catch (Exception e) {
            e.printStackTrace();
        }

        sourceBundlePath = go.warpify.Warpify.sourcePath();
    }

    @Override
    public List<NativeModule> createNativeModules(ReactApplicationContext reactContext) {
        List<NativeModule> modules = new ArrayList<>();
        modules.add(new WarpifyModule(reactContext));
        return modules;
    }

    @Override
    public List<Class<? extends JavaScriptModule>> createJSModules() {
        return Collections.emptyList();
    }

    @Override
    public List<ViewManager> createViewManagers(ReactApplicationContext reactContext) {
        return Collections.emptyList();
    }

    class WarpifyModule extends ReactContextBaseJavaModule {
        public WarpifyModule(final ReactApplicationContext reactContext) {
            super(reactContext);

            // we need to pass the internal callback. the internal callback is used to know whether
            // we have a new update or not
            Warpify.setReload(new go.warpify.Callback() {
                @Override
                public void do_(long kind, String path) {
                    sourceBundlePath = path;
                    reloadBundle();
                }
            });
        }

        private void reloadBundle() {
            final Activity currentActivity = getCurrentActivity();
            if (currentActivity != null) {
                Intent intent = currentActivity.getIntent();
                currentActivity.finish();
                currentActivity.startActivity(intent);
                // I don't really know if we need this.
                currentActivity.recreate();
            }
        }

        @ReactMethod
        public void cycles(Callback errorCallback, Callback successCallback) {
            try {
                String result = Warpify.cycles();
                successCallback.invoke(result);
            } catch (Exception e) {
                errorCallback.invoke(e.getMessage());
            }
        }

        @ReactMethod
        public void remoteVersions(int cycleId, Callback errorCallback, Callback successCallback) {
            try {
                String result = Warpify.remoteVersions(cycleId);
                successCallback.invoke(result);
            } catch (Exception e) {
                errorCallback.invoke(e.getMessage());
            }
        }

        @ReactMethod
        public void localVersions(int cycleId, Callback errorCallback, Callback successCallback) {
            try {
                String result = Warpify.localVersions(cycleId);
                successCallback.invoke(result);
            } catch (Exception e) {
                errorCallback.invoke(e.getMessage());
            }
        }

        @ReactMethod
        public void latestVersion(int cycleId, Callback errorCallback, Callback successCallback) {
            try {
                String result = Warpify.latest(cycleId);
                successCallback.invoke(result);
            } catch (Exception e) {
                errorCallback.invoke(e.getMessage());
            }
        }

        @ReactMethod
        public void downloadVersion(int cycleId, String version, Callback errorCallback, Callback successCallback) {
            try {
                Warpify.downloadVersion(cycleId, version);
                successCallback.invoke();
            } catch (Exception e) {
                errorCallback.invoke(e.getMessage());
            }
        }

        @ReactMethod
        public void reload(int cycleId, String version, Callback errorCallback, Callback successCallback) {
            try {
                Warpify.reload(cycleId, version);
                successCallback.invoke();
            } catch (Exception e) {
                errorCallback.invoke(e.getMessage());
            }
        }

        @Override
        public String getName() {
            return "WarpifyManager";
        }
    }
}