package com.pressly.warpdrive;

import android.content.Context;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;

import com.facebook.react.ReactPackage;
import com.facebook.react.bridge.JavaScriptModule;
import com.facebook.react.bridge.NativeModule;
import com.facebook.react.bridge.ReactApplicationContext;
import com.facebook.react.bridge.ReactContextBaseJavaModule;
import com.facebook.react.uimanager.ViewManager;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import warpdrive.Warpdrive;

public class WarpdrivePackage implements ReactPackage {

    public WarpdrivePackage(Context context, String app, String rolloutAt, String serverAddr, String deviceCert, String deviceKey, String caCert) {
        String documentPath = context.getFilesDir().getAbsolutePath();
        String bundlePath = getBundlePath(context, documentPath, new String[]{ deviceCert, deviceKey, caCert });
        String platform = "android";
        String bundleVersion = getBundleVersion(context);

        deviceCert = joinPath(bundlePath, deviceCert);
        deviceKey = joinPath(bundlePath, deviceKey);
        caCert = joinPath(bundlePath, caCert);

        try {
            Warpdrive.init(bundlePath, documentPath, platform, app, rolloutAt, bundleVersion, serverAddr, deviceCert, deviceKey, caCert);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public static String sourceBundle() {
        String path = Warpdrive.bundlePath();
        if (path == null || path.equals("")) {
            path = null;
        }
        return path;
    }

    @Override
    public List<NativeModule> createNativeModules(ReactApplicationContext reactContext) {
        List<NativeModule> modules = new ArrayList<>();
        modules.add(new WarpdriveModule(reactContext));
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

    // since we can't get the bundle path in android,
    // we need to copy all certificate files into document folder under `warpdrive_bundle`
    // this folder will keep overridden every time app boots up
    private static String getBundlePath(Context context, String documentPath, String[] filenames) {
        String bundlePath = joinPath(documentPath, "warpdrive_assets");

        // create a warpdrive_assets folder
        new File(bundlePath).mkdirs();

        InputStream src = null;
        OutputStream dst = null;

        try {
            for (String filename : filenames) {
                src = context.getAssets().open(filename);
                dst = new FileOutputStream(joinPath(bundlePath, filename));
                if (copy(dst, src) == 0) {
                    System.out.printf("couldn't copy %s", filename);
                    break;
                }
                src.close();
                dst.close();
            }
        } catch (IOException e) {
            if (src != null) {
                try {
                    src.close();
                } catch (IOException e1) {
                    e1.printStackTrace();
                }
            }
        }

        return bundlePath;
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

    private static String joinPath(String a, String b) {
        return new File(a, b).getPath();
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


    class WarpdriveModule extends ReactContextBaseJavaModule {
        public WarpdriveModule(final ReactApplicationContext reactContext) {
            super(reactContext);
        }

        @Override
        public String getName() {
            return "WarpdriveManager";
        }
    }
}
