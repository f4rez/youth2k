//
//  ViewController.swift
//  Youth
//
//  Created by Josef Svensson on 2017-07-08.
//  Copyright © 2017 Josef Svensson. All rights reserved.
//

import UIKit

class ViewController: UIViewController, UIWebViewDelegate {
    @IBOutlet var schedule: UIImageView!
    @IBOutlet var downloads: UIImageView!
    @IBOutlet var speakers: UIImageView!
    
    @IBOutlet var splash: UIImageView!
    @IBOutlet var mWebView: UIWebView!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view, typically from a nib.
        mWebView.delegate = self
        mWebView.scrollView.bounces = false
        mWebView.loadRequest(URLRequest(url: URL(string: "http://youth.livetsord.se")!))
        
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }

    func webViewDidFinishLoad(_ webView: UIWebView) {
        print("Finished")
        UIView.animate(withDuration: 0.9, animations: {
            
            self.splash.alpha = 0.0
        })
    }
    

}

