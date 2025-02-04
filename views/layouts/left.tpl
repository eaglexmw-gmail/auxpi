<div id="menu" class="mdui-drawer mdui-drawer-close">
    <div class="mdui-list" mdui-collapse="{accordion: true}">
        {{range $elem := .Stores}}
            <a class="mdui-list-item mdui-ripple" href="{{$elem.Router}}">
                <i class="mdui-list-item-icon mdui-icon material-icons iconfont icon-{{$elem.Icon}} mdui-text-color-{{$elem.Color}}"></i>
                <div class="mdui-list-item-content">{{$elem.Name}}</div>
            </a>
        {{end}}
        <script>function pages_ctrl_about_show() {location.href = location.origin + window.sub_folder + "{{urlfor "PagesController.AboutShow"}}"}</script>
        <a class="mdui-list-item mdui-ripple " href="javascript:void(0);" onclick="pages_ctrl_about_show()">
            <i class="mdui-list-item-icon mdui-icon material-icons iconfont icon-about mdui-text-color-teal"></i>
            <div class="mdui-list-item-content">说明</div>
        </a>
        <a class="mdui-list-item mdui-ripple " href="https://github.com/eaglexmw-gmail/AUXPI" target="_blank">
            <i class="mdui-list-item-icon mdui-icon material-icons iconfont icon-icons-code mdui-text-color-red"></i>
            <div class="mdui-list-item-content">项目地址</div>
        </a>
    </div>
</div>