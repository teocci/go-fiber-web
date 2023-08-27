/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-6월-10
 */
import PositionModule from './modules/position-module.js'
import LoaderComponent from './components/elements/loader-component.js'

window.onload = () => {
    console.log('init')

    loaderComponent = new LoaderComponent()
    mainModule = PositionModule.instance
}