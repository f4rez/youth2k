//
//  ViewController.swift
//  Youth
//
//  Created by Josef Svensson on 2017-07-08.
//  Copyright © 2017 Josef Svensson. All rights reserved.
//

import UIKit
import Photos

class ViewController: UIViewController, UIWebViewDelegate {
    
    
    @IBOutlet var splash: UIImageView!
    @IBOutlet var mWebView: UIWebView!
    @IBOutlet var mWebView2: UIWebView!
    
    @IBOutlet var home: UIImageView!
    @IBOutlet var speakers: UIImageView!
    @IBOutlet var downloads: UIImageView!
    @IBOutlet var schedule: UIImageView!
    @IBOutlet var spinner: UIActivityIndicatorView!
    @IBOutlet var progress: UIProgressView!
    @IBOutlet var bottombar: UIImageView!
    @IBOutlet var fade: UIImageView!
    
    var counterWebview = 0
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view, typically from a nib.
        mWebView.delegate = self
        mWebView.scrollView.bounces = false
        mWebView2.delegate = self
        mWebView2.scrollView.bounces = false
        loadRequestWebview(mWebView: mWebView, mUrl:"http://youthapp.livetsord.se/hem" )
        schedule.isUserInteractionEnabled = true
        downloads.isUserInteractionEnabled = true
        speakers.isUserInteractionEnabled = true
        home.isUserInteractionEnabled = true
        home.isHighlighted = true
        
        let hTapGestureRec = UITapGestureRecognizer(target: self, action: #selector(homeTapped(tapGestureRecognizer:)))
        let tapGestureRecognizer = UITapGestureRecognizer(target: self, action: #selector(scheduleTapped(tapGestureRecognizer:)))
        let sPtapGestureRecognizer = UITapGestureRecognizer(target: self, action: #selector(sPTapped(tapGestureRecognizer:)))
        let dOtapGestureRecognizer = UITapGestureRecognizer(target: self, action: #selector(dOTapped(tapGestureRecognizer:)))
        spinner.startAnimating()
        
        schedule.addGestureRecognizer(tapGestureRecognizer)
        downloads.addGestureRecognizer(dOtapGestureRecognizer)
        speakers.addGestureRecognizer(sPtapGestureRecognizer)
        home.addGestureRecognizer(hTapGestureRec)
        
    }
    
    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    func webViewDidFinishLoad(_ webView: UIWebView) {
        print("Finished")
        self.spinner.stopAnimating()
        self.progress.setProgress(1.0, animated: true)
        Timer.scheduledTimer(timeInterval: 1.0, target: self, selector: #selector(self.progressDone(sender:)), userInfo: nil, repeats: false)
        
        UIView.animate(withDuration: 0.9, animations: {
            self.spinner.alpha = 0.0
            self.splash.alpha = 0.0
        })
        
        if counterWebview % 2 == 0 {
            if let remove = webView.request?.url?.absoluteString.hasSuffix("/hem/") {
                UIView.animate(withDuration: 0.9, animations: {
                    self.mWebView2.alpha = 0
                    self.mWebView.alpha = 1
                    if remove {
                        self.fade.alpha = 0
                        self.bottombar.alpha = 0
                    } else {
                        self.fade.alpha = 1.0
                        self.bottombar.alpha = 1.0
                    }
                })
            }
            
        } else {
            if let remove = webView.request?.url?.absoluteString.hasSuffix("/hem/") {
                
                UIView.animate(withDuration: 0.9, animations: {
                    self.mWebView.alpha = 0
                    self.mWebView2.alpha = 1
                    if remove {
                        self.fade.alpha = 0
                        self.bottombar.alpha = 0
                    } else {
                        self.fade.alpha = 1.0
                        self.bottombar.alpha = 1.0
                    }
                })
            }
            
        }
        counterWebview += 1
    }
    
    func webView(_ webView: UIWebView, shouldStartLoadWith request: URLRequest, navigationType: UIWebViewNavigationType) -> Bool {
        print("Entered should load")
        if let mUrl = request.url?.absoluteString {
            if mUrl.hasSuffix(".png") || mUrl.hasSuffix(".jpg") {
                let url = URL(string: mUrl)
                let data = try? Data(contentsOf: url!)
                if let img = UIImage(data: data!) {
                    UIImageWriteToSavedPhotosAlbum(img, nil, nil, nil);
                    let alert = UIAlertController(title: "Bilden är sparad i din kamerarulle", message: "", preferredStyle: UIAlertControllerStyle.alert)
                    alert.addAction(UIAlertAction(title: "Ok", style: UIAlertActionStyle.default, handler: nil))
                    self.present(alert, animated: true, completion: nil)
                } else {
                    print("piss")
                }
                return false
            }
        }
        print("Return true, should load")
        return true
    }
    
    func loadRequestWebview(mWebView: UIWebView, mUrl: String) {
        if counterWebview % 2 == 0 {
            mWebView.loadRequest(URLRequest(url: URL(string: mUrl)!,cachePolicy: URLRequest.CachePolicy.returnCacheDataElseLoad, timeoutInterval: 20.0))
        } else {
            mWebView2.loadRequest(URLRequest(url: URL(string: mUrl)!,cachePolicy: URLRequest.CachePolicy.returnCacheDataElseLoad, timeoutInterval: 20.0))
        }
        
    }
    
    
    func unHighlightAll() {
        self.speakers.isHighlighted = false
        self.downloads.isHighlighted = false
        self.schedule.isHighlighted = false
        self.home.isHighlighted = false
    }
    
    func homeTapped(tapGestureRecognizer: UITapGestureRecognizer)
    {
        let tappedImage = tapGestureRecognizer.view as! UIImageView
        if !tappedImage.isHighlighted {
            unHighlightAll()
            startProgressBar()
            tappedImage.isHighlighted = true
            loadRequestWebview(mWebView: self.mWebView, mUrl:"http://youthapp.livetsord.se/hem")
        }
    }
    
    
    func scheduleTapped(tapGestureRecognizer: UITapGestureRecognizer)
    {
        let tappedImage = tapGestureRecognizer.view as! UIImageView
        if !tappedImage.isHighlighted {
            unHighlightAll()
            startProgressBar()
            tappedImage.isHighlighted = true
            loadRequestWebview(mWebView: self.mWebView, mUrl:"http://youthapp.livetsord.se/schema")
        }
    }
    
    func sPTapped(tapGestureRecognizer: UITapGestureRecognizer)
    {
        let tappedImage = tapGestureRecognizer.view as! UIImageView
        if !tappedImage.isHighlighted {
            unHighlightAll()
            startProgressBar()
            tappedImage.isHighlighted = true
            loadRequestWebview(mWebView: self.mWebView, mUrl:"http://youthapp.livetsord.se/talare")
        }
    }
    func dOTapped(tapGestureRecognizer: UITapGestureRecognizer)
    {
        let tappedImage = tapGestureRecognizer.view as! UIImageView
        if !tappedImage.isHighlighted {
            unHighlightAll()
            startProgressBar()
            tappedImage.isHighlighted = true
            loadRequestWebview(mWebView: self.mWebView, mUrl:"http://youthapp.livetsord.se/downloads")
        }
        
    }
    
    func openLinkFromNotif(link: String) {
        loadRequestWebview(mWebView: self.mWebView, mUrl: link)
    }
    
    func startProgressBar() {
        self.progress.setProgress(0.0, animated: false)
        self.progress.isHidden = false
        Timer.scheduledTimer(timeInterval: 1/60, target: self, selector: #selector(self.updateTimer(sender:)), userInfo: nil, repeats: true)
    }
    
    func updateTimer(sender: Timer){
        if self.progress.progress < 0.9 {
            self.progress.setProgress(self.progress.progress+0.001, animated: true)
        } else {
            sender.invalidate()
        }
        
    }
    func progressDone(sender: Timer){
        self.progress.isHidden = true
        self.progress.setProgress(0, animated: false)
    }
    
    func removeBottomBar() {
        UIView.animate(withDuration: 1.9, animations: {
            self.fade.alpha = 0.0
            self.bottombar.alpha = 0.0
        })
    }
    
    func reviewBottomBar() {
        UIView.animate(withDuration: 0.2, animations: {
            self.fade.alpha = 1.0
            self.bottombar.alpha = 1.0
        })
    }
}


