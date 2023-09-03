/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-6월-10
 */
import LoaderComponent from './components/elements/loader-component.js'
import AdsModule from './modules/ads-module.js'

window.onload = () => {
    console.log('init')

    loaderComponent = new LoaderComponent()
    mainModule = AdsModule.instance
}