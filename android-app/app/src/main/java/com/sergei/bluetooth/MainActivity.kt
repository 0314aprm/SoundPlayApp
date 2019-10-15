package com.sergei.bluetooth

import android.os.Bundle
import android.support.design.widget.FloatingActionButton
import android.support.design.widget.Snackbar
import android.support.v7.app.AppCompatActivity
import android.support.v7.widget.Toolbar
import android.view.View
import android.view.Menu
import android.view.MenuItem
import com.sergei.bluetooth.bluetooth.BluetoothConnectionService

class MainActivity : AppCompatActivity() {

    private val bluetoothConnectionService : BluetoothConnectionService = BluetoothConnectionService(this)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
    }
}
