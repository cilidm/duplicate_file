package util

import (
	"fmt"
	"github.com/cilidm/toolbox/file"
	"github.com/cilidm/toolbox/logging"
	"os"
	"path"
)

var (
	templates     = []string{"start.html", "home.html", "all.html", "pic.html", "video.html", "file_manager.html"}
	statics       = []string{"static/layui_exts/fileManager/fileManager.js"}
	HomeTemplate  = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>DuplicateFiles</title><meta name="renderer" content="webkit"><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"><link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css"></head><body style="width:98%;margin:0 auto"><table class="layui-hide" id="test" lay-filter="test"></table><script type="text/html" id="barDemo"><a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del">删除</a></script><script src="https://www.layuicdn.com/layui/layui.js"></script><script>layui.use(["table","jquery"],function(){var e=layui.table,l=layui.jquery;e.render({elem:"#test",url:"/json",toolbar:"#toolbarDemo",defaultToolbar:["filter","exports","print"],title:"DuplicateFiles",cols:[[{field:"name",title:"文件名称"},{field:"dir",title:"文件夹",minWidth:450},{field:"size",title:"文件大小",sort:!0},{field:"mod_time",title:"修改日期",width:180},{field:"md",title:"MD5",sort:!0,width:275},{fixed:"right",title:"操作",toolbar:"#barDemo",width:90}]],page:!0,even:!0,skin:"line",limits:[20,50,100,500,1e3],limit:20}),e.on("tool(test)",function(e){var i=e.data;"del"===e.event&&layer.confirm("是否删除文件"+i.name,function(t){l.post("/del",{path:i.path},function(e){layer.close(t),200==e.code?layer.msg("操作成功",{icon:1,shade:.3,time:1e3},function(){location.reload()}):layer.msg(e.msg)})})})})</script></body></html>`
	StartTemplate = `<!DOCTYPE html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"><title>DuplicateFile</title><link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css"><style>.layui-header{height:50px}.layui-layout-admin .layui-logo{line-height:50px;height:50px;font-size:18px}.layui-nav layui-layout-right{line-height:50px;height:50px}.layui-layout-admin .layui-side{top:50px}.layui-layout-admin .layui-body{top:50px}.pp-nav-item{line-height:50px!important;height:50px!important}.pp-nav-child a{line-height:30px!important;height:30px!important}.pp-nav-item a{line-height:50px;height:50px}.layui-nav .layui-nav-item .layui-nav-itemed{line-height:50px;height:50px}.layui-nav-child{top:50px}.layui-layout-admin .layui-footer{line-height:35px;height:35px}.pp-side-fold{height:30px;background-color:#4a5064;color:#aeb9c2;line-height:30px;text-align:center;cursor:pointer}.back_space1{margin-right:5px}a.pointer{cursor:pointer}.pp-ddsided{width:50px}a.pp-pointer{padding-right:0!important}.pp-sided{width:50px}.pp-main{left:50px!important}.pp-tab{margin-top:0}.pp-tab-title{background:#000}.layui-this{background:#fff}.layui-tab-title li.pp-tab-li{border:0}</style></head><body class="layui-layout-body"><div class="layui-layout layui-layout-admin"><div class="layui-header"><div class="layui-logo">&nbsp;&nbsp;&nbsp;&nbsp;DuplicateFile<font style="font-size:12px"> version1.0</font></div><ul class="layui-nav layui-layout-left"></ul><ul class="layui-nav layui-layout-right pp-nav-item"><li class="layui-nav-item"><a href="http://github.com/cilidm" target="_blank">GitHub</a></li></ul></div><div class="layui-side layui-bg-black pp-side"><div class="layui-side-scroll" style="width:100%"><ul class="layui-nav layui-nav-tree" lay-filter="nav-side"><li class="layui-nav-item layui-nav-itemed layui-this"><a data-url="/main" data-icon="fa-home" data-title="系统首页" data-id="0" class="pointer"><i class="fa fa-home back_space1"></i><span>重复文件</span></a></li><li class="layui-nav-item"><a class="" href="javascript:;"><i class="fa fa-home"></i> &nbsp;&nbsp;<span>查看文件</span></a><dl class="layui-nav-child pp-nav-childs"><dd><a data-url="/index" data-icon="fa fa-home" data-title="全部文件" class="pointer" data-id="1"><i class="fa fa-home"></i> &nbsp;&nbsp;<span>全部文件</span></a></dd><dd><a data-url="/pic" data-icon="fa fa-home" data-title="图片文件" class="pointer" data-id="2"><i class="fa fa-home"></i> &nbsp;&nbsp;<span>图片文件</span></a></dd><dd><a data-url="/video" data-icon="fa fa-home" data-title="视频文件" class="pointer" data-id="3"><i class="fa fa-home"></i> &nbsp;&nbsp;<span>视频文件</span></a></dd><dd><a data-url="/file" data-icon="fa fa-home" data-title="文件展示" class="pointer" data-id="4"><i class="fa fa-home"></i> &nbsp;&nbsp;<span>文件展示</span></a></dd></dl></li></ul></div></div><div class="layui-body" style="overflow:hidden"><div class="layui-tab pp-tab" lay-filter="main_tab" lay-allowClose="true"><ul class="layui-tab-title" style="background:#efefef"><li class="pp-tab-li layui-this" id="default_tab" lay-id="0"><i class="fa fa-home back_space1"></i>重复文件</li></ul><div class="layui-tab-content" style="padding:0"><div class="layui-tab-item layui-show" style="margin:0;overflow:hidden"><iframe src="/main" frameborder="0" scrolling="yes"></iframe></div></div></div></div><div class="layui-footer" style="border-top:2px solid #e4e4e4">© Power by <a href="http://www.github.com/" target="_blank">.version</a></div></div><script src="https://www.layuicdn.com/layui/layui.js"></script><script>var $,element,width,height;function delHtmlTag(i){return i.replace(/<[^>]+>/g,"")}function getDelimiterLastString(i,t){return arr=i.split(t),1<arr.length?arr[arr.length-1]:i}function deleteCurrentTab(){$(parent.document).find("ul.layui-tab-title").children("li.layui-this").find(".layui-tab-close").click();var i=$(parent.document).find("div.layui-tab-content").find("div.layui-show").find("iframe").contents();$(i[0]).find("body").find("#reload").click()}function getCheckboxValue(i){var t=new Array;return $("input:checkbox[name="+i+"]:checked").each(function(){console.log($(this).val()),t.push($(this).val())}),t.join(",")}function openTab(i,t,a,e){var n,d;i&&t&&a&&(e=e||" fa-clock-o ",n=0,$(".layui-tab-title").find("li").each(function(){$(this).attr("lay-id")==a&&(n=1)}),1==n||(d='<iframe src="'+i+'" scrolling="yes" width="'+width+'" height="'+height+'" frameborder="0"></iframe>',t='<i class="fa '+e+' back_space1"></i>'+t,element.tabAdd("main_tab",{title:t,content:d,id:a})),element.tabChange("main_tab",a))}layui.use(["element","jquery","layer"],function(){element=layui.element,$=layui.jquery;var i=layui.layer;function t(){height=$(".layui-body").height()-40,width=$(".layui-body").width(),$(".ayui-tab-content").width(width),$(".ayui-tab-content").height(height),$(".layui-tab-item").height(height),$(".layui-tab-item").find("iframe").height(height),$(".layui-tab-item").find("iframe").width(width)}t(),window.onresize=function(){t()},$(".pp-side-fold").on("click",function(){50<$(".layui-side").width()?($(".layui-side").width(50),$(this).parent().find("span").hide(),$(".layui-body").addClass("pp-main"),$(".layui-footer").addClass("pp-main"),$(".layui-nav-child").find("dd").addClass("pp-ddsided"),$(".layui-nav-child").find("a").addClass("pp-pointer")):($(".layui-side").width(200),$(this).parent().find("span").show(),$(".layui-body").removeClass("pp-main"),$(".layui-footer").removeClass("pp-main"),$(".layui-nav-child").find("dd").removeClass("pp-ddsided"),$(".layui-nav-child").find("a").removeClass("pp-pointer")),t()}),element.on("nav(nav-side)",function(i){var t,a,e=i.attr("data-url"),n=i.attr("data-title"),d=i.attr("data-id"),l=i.attr("data-icon");e&&n&&d&&(t=0,$(".layui-tab-title").find("li").each(function(){$(this).attr("lay-id")==d&&(t=1)}),1==t||(a='<iframe src="'+e+'" scrolling="yes" width="'+width+'" height="'+height+'" frameborder="0"></iframe>',n='<i class="fa '+l+' back_space1"></i>'+n,element.tabAdd("main_tab",{title:n,content:a,id:d})),element.tabChange("main_tab",d))}),$(".layui-tab-title").on("click","li",function(){var i=$(this).attr("lay-id");i&&$(".layui-nav-item").find("a").each(function(){$(this).attr("data-id")==i&&$(this).click()})}),$(".pp-nav-childs").find("a").hover(function(){i.tips($(this).attr("data-title"),$(this),{time:1e3})}),element.render()})</script></body></html>`
	AllTemplate   = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>DuplicateFiles</title><meta name="renderer" content="webkit"><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"><link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css"></head><body style="width:98%;margin:0 auto"><table class="layui-hide" id="test" lay-filter="test"></table><script type="text/html" id="barDemo">{{/* <a class="layui-btn layui-btn-xs" lay-event="edit">查看</a>*/}}<a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del">删除</a></script><script src="https://www.layuicdn.com/layui/layui.js"></script><script>layui.use(["table","jquery"],function(){var e=layui.table,l=layui.jquery;e.render({elem:"#test",url:"/index_json",toolbar:"#toolbarDemo",defaultToolbar:["filter","exports","print"],title:"DuplicateFiles",cols:[[{field:"name",title:"文件名称"},{field:"dir",title:"文件夹",minWidth:450},{field:"size",title:"文件大小",width:140,sort:!0},{field:"mod_time",title:"修改日期"},{field:"md",title:"MD5",sort:!0,width:275},{fixed:"right",title:"操作",toolbar:"#barDemo",width:150}]],page:!0,even:!0,skin:"line",limits:[20,50,100,500,1e3],limit:20}),e.on("tool(test)",function(e){var i=e.data;"del"===e.event&&layer.confirm("是否删除文件"+i.name,function(t){l.post("/del",{path:i.path},function(e){layer.close(t),200==e.code?layer.msg("操作成功",{icon:1,shade:.3,time:1e3},function(){location.reload()}):layer.msg(e.msg)})})})})</script></body></html>`
	PicTemplate   = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>DuplicateFiles</title><meta name="renderer" content="webkit"><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"><link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css"></head><body style="width:98%;margin:0 auto"><table class="layui-hide" id="test" lay-filter="test"></table><script type="text/html" id="barDemo"><a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del">删除</a></script><script src="https://www.layuicdn.com/layui/layui.js"></script><script>layui.use(["table","jquery"],function(){var e=layui.table,l=layui.jquery;e.render({elem:"#test",url:"/pic_json",toolbar:"#toolbarDemo",defaultToolbar:["filter","exports","print"],title:"DuplicateFiles",cols:[[{field:"name",title:"文件名称"},{field:"dir",title:"文件夹",minWidth:450},{field:"size",title:"文件大小",width:140,sort:!0},{field:"mod_time",title:"修改日期"},{field:"md",title:"MD5",sort:!0,width:275},{fixed:"right",title:"操作",toolbar:"#barDemo",width:150}]],page:!0,even:!0,skin:"line",limits:[20,50,100,500,1e3],limit:20}),e.on("tool(test)",function(e){var i=e.data;"del"===e.event&&layer.confirm("是否删除文件"+i.name,function(t){l.post("/del",{path:i.path},function(e){layer.close(t),200==e.code?layer.msg("操作成功",{icon:1,shade:.3,time:1e3},function(){location.reload()}):layer.msg(e.msg)})})})})</script></body></html>`
	VideoTemplate = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>DuplicateFiles</title><meta name="renderer" content="webkit"><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"><link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css"></head><body style="width:98%;margin:0 auto"><table class="layui-hide" id="test" lay-filter="test"></table><script type="text/html" id="barDemo"><a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del">删除</a></script><script src="https://www.layuicdn.com/layui/layui.js"></script><script>layui.use(["table","jquery"],function(){var e=layui.table,l=layui.jquery;e.render({elem:"#test",url:"/video_json",toolbar:"#toolbarDemo",defaultToolbar:["filter","exports","print"],title:"DuplicateFiles",cols:[[{field:"name",title:"文件名称"},{field:"dir",title:"文件夹",minWidth:450},{field:"size",title:"文件大小",width:140,sort:!0},{field:"mod_time",title:"修改日期"},{field:"md",title:"MD5",sort:!0,width:275},{fixed:"right",title:"操作",toolbar:"#barDemo",width:150}]],page:!0,even:!0,skin:"line",limits:[20,50,100,500,1e3],limit:20}),e.on("tool(test)",function(e){var i=e.data;"del"===e.event&&layer.confirm("是否删除文件"+i.name,function(t){l.post("/del",{path:i.path},function(e){layer.close(t),200==e.code?layer.msg("操作成功",{icon:1,shade:.3,time:1e3},function(){location.reload()}):layer.msg(e.msg)})})})})</script></body></html>`
	FileManager   = `"use strict";var _typeof="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"==typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t};layui.define(["jquery","layer","laypage"],function(t){var e=layui.jquery,i=layui.layer,a=layui.laypage,n={config:{test:"test",thumb:{nopic:"",width:100,height:100},icon_url:"ico/",btn_upload:!0,btn_create:!0},cache:{},index:layui.fm?layui.fm.index+1e4:0,set:function(t){var i=this;return i.config=e.extend({},i.config,t),i},on:function(t,e){return layui.onevent.call(this,"fileManager",t,e)},dirRoot:[{path:"",name:"根目录"}],v:"1.0.1.2019.12.26"},o=function t(){var e=this,i=e.config,a=i.id||i.index;return a&&(t.that[a]=e,t.config[a]=i),{config:i,reload:function(t){e.reload.call(e,t)}}},l=function(t){var e=o.config[t];return e||hint.error("The ID option was not found in the fm instance"),e||null},r=function(t){var i=this;i.config=e.extend({},i.config,n.config,t),o.that={},o.config={},i.render()};r.prototype.render=function(){var t=this,i=t.config;if(i.elem=e(i.elem),i.where=i.where||{},i.id=i.id||i.elem.attr("id")||t.index,i.request=e.extend({pageName:"page",limitName:"limit"},i.request),i.response=e.extend({statusName:"code",statusCode:0,msgName:"msg",dataName:"data",countName:"count"},i.response),"object"===_typeof(i.page)&&(i.limit=i.page.limit||i.limit,i.limits=i.page.limits||i.limits,t.page=i.page.curr=i.page.curr||1,delete i.page.elem,delete i.page.jump),!i.elem[0])return t;var a="";i.btn_create&&(a+='<button type="button" class="layui-btn layui-btn-primary layui-btn-sm" id="new_dir">建文件夹</button>'),i.btn_upload&&(a+='<button type="button" class="layui-btn layui-btn-primary layui-btn-sm" id="uploadfile">上传文件</button>');var n='<div class="layui-card" ><div class="layui-card-body"><div class="layui-btn-group tool_bar">'+a+'<button type="button" class="layui-btn layui-btn-primary layui-btn-sm" id="back"><i class="layui-icon layui-icon-left line"></i></button></div><div class="layui-inline path_bar" id=""><a ><i class="layui-icon layui-icon-more-vertical line" ></i>根目录</a></div></div><hr><div class="layui-card-body"><div class="file-body layui-form" style=""><ul class="fileManager layui-row fm_body layui-col-space10" ></ul></div><hr><div ><div class="layui_page_'+i.id+'" id="layui_page_'+i.id+'"></div></div></div>';i.elem.html(n),i.index=t.index,t.key=i.id||i.index,t.layPage=i.elem.find(".layui_page_"+i.id),t.layBody=i.elem.find(".fm_body"),t.layPathBar=i.elem.find(".path_bar"),t.layToolBar=i.elem.find(".tool_bar"),t.pullData(t.page),t.events()},r.prototype.page=1,r.prototype.pullData=function(t){var i=this,a=i.config,n=a.request,o=a.response,l=!1;if(i.startTime=(new Date).getTime(),a.url){var r={};r[n.pageName]=t,r[n.limitName]=a.limit;var c=e.extend(r,a.where);a.contentType&&0==a.contentType.indexOf("application/json")&&(c=JSON.stringify(c)),i.loading(),e.ajax({type:a.method||"get",url:a.url,contentType:a.contentType,data:c,async:!1,dataType:"json",headers:a.headers||{},success:function(e){"function"==typeof a.parseData&&(e=a.parseData(e)||e),e[o.statusName]!=o.statusCode?i.errorView(e[o.msgName]||'返回的数据不符合规范，正确的成功状态码应为："'+o.statusName+'": '+o.statusCode):(i.renderData(e,t,e[o.countName]),a.time=(new Date).getTime()-i.startTime+" ms"),"function"==typeof a.done&&a.done(e,t,e[o.countName]),l=!0},error:function(t,e){i.errorView("数据接口请求异常："+e)}})}return l},r.prototype.renderData=function(t,i,o){var l=this,r=l.config,c=t[r.response.dataName]||[],u="";if(layui.each(c,function(t,e){var i=void 0,a=void 0;switch(a=e.type,e.type){case"directory":i='<div  style="width:'+r.thumb.width+"px;height:"+r.thumb.height+"px;line-height:"+r.thumb.height+'px"><img src="ico/dir.png" style="vertical-align:middle;"></div>',a="DIR";break;default:i="png"==e.type||"gif"==e.type||"jpg"==e.type||"image"==e.type?'<img src="'+e.thumb+'" width="'+r.thumb.width+'" height="'+r.thumb.height+'" onerror=\'this.src="'+r.thumb.nopic+"\"'  />":'<div  style="width:'+r.thumb.width+"px;height:"+r.thumb.height+"px;line-height:"+r.thumb.height+'px"><img src="'+r.icon_url+e.type+'.png"  onerror=\'this.src="'+r.thumb.nopic+"\"' /></div>"}u+='<li style="display:inline-block" data-type="'+a+'" data-index="'+t+'"><div class="content" align="center">'+i+'<p class="layui-elip" title="'+e.name+'">'+e.name+" </p></div></li>"}),r.elem.find(".fileManager").html(u),n.cache[r.id]=c,l.layPage[0==o||0===c.length&&1==i?"addClass":"removeClass"]("layui-hide"),0===c.length)return l.errorView("空目录");r.page&&(r.page=e.extend({elem:"layui_page_"+r.id,count:o,limit:r.limit,limits:r.limits||[10,20,30,40,50,60,70,80,90],groups:3,layout:["prev","page","next","skip","count","limit"],prev:'<i class="layui-icon">&#xe603;</i>',next:'<i class="layui-icon">&#xe602;</i>',jump:function(t,e){e||(l.page=t.curr,r.limit=t.limit,l.pullData(t.curr))}},r.page),r.page.count=o,a.render(r.page))},r.prototype.updatePathBar=function(){var t=this,e=t.config,i=n.dirRoot[n.dirRoot.length-1];e.where={path:i.path},0!=t.pullData(1)&&(t.layPathBar.html(""),n.dirRoot.map(function(e,i,a){var n=0==i?"layui-icon-more-vertical":"layui-icon-right",o='<i class="layui-icon '+n+'"></i><a  data-path="'+e.path+'" data-name="'+e.name+'" >'+e.name+"</a>";t.layPathBar.append(o)}))},r.prototype.events=function(){var t=this,a=t.config,o=(e("body"),a.elem.attr("lay-filter"));t.layBody.on("click","li",function(){l.call(this,"pic")}),t.layBody.on("click","li[data-type=DIR]",function(){var i=e(this),o=n.cache[a.id];o=o[i.data("index")]||{},n.dirRoot.push({path:o.path,name:o.name}),t.updatePathBar()}),t.layToolBar.on("click","#back",function(){e(this);if(1==n.dirRoot.length)return i.msg("已经是根目录");n.dirRoot.length>1&&n.dirRoot.pop(),t.updatePathBar()}),t.layToolBar.on("click","#uploadfile",function(){var t=e(this);layui.event.call(this,"fileManager","uploadfile("+o+")",{obj:t,path:n.dirRoot[n.dirRoot.length-1].path})}),t.layToolBar.on("click","#new_dir",function(){var t=e(this);i.prompt({title:"请输入新文件夹名字",formType:0},function(e,a){i.close(a),layui.event.call(this,"fileManager","new_dir("+o+")",{obj:t,folder:e,path:n.dirRoot[n.dirRoot.length-1].path})})});var l=function(t){var i=e(this),l=n.cache[a.id],r=i.data("index");"DIR"!=i.data("type")&&(l=l[r]||{},layui.event.call(this,"fileManager",t+"("+o+")",{obj:i,data:l}))}},r.prototype.loading=function(t){var i=this;i.config.loading&&(t?(i.layInit&&i.layInit.remove(),delete i.layInit,i.layBox.find(ELEM_INIT).remove()):(i.layInit=e(['<div class="layui-table-init">','<i class="layui-icon layui-icon-loading layui-anim layui-anim-rotate layui-anim-loop"></i>',"</div>"].join("")),i.layBox.append(i.layInit)))},r.prototype.errorView=function(t){i.msg(t)},r.prototype.reload=function(t){var i=this;t=t||{},delete i.haveInit,t.data&&t.data.constructor===Array&&delete i.config.data,i.config=e.extend(!0,{},i.config,t),i.render()},n.reload=function(t,e){if(l(t)){var i=o.that[t];return i.reload(e),o.call(i)}},n.render=function(t){var e=new r(t);return o.call(e)},t("fileManager",n)});`
)

func getTemInfo(fileName string) string {
	var res = make(map[string]string)
	res["start.html"] = StartTemplate
	res["home.html"] = HomeTemplate
	res["all.html"] = AllTemplate
	res["pic.html"] = PicTemplate
	res["video.html"] = VideoTemplate
	res["file_manager.html"] = FileTemplate
	return res[fileName]
}

func getStatics(staticName string) string {
	var res = make(map[string]string)
	res["static/layui_exts/fileManager/fileManager.js"] = FileManager
	return res[staticName]
}

func CheckTemplate() {
	for _, v := range templates {
		temFile := "views/" + v
		if CheckFileIsExist(temFile) == false {
			file.IsNotExistMkDir("views")
			f, err := os.Create(temFile)
			if err != nil {
				fmt.Println(err.Error())
				logging.Fatal(err.Error())
			}
			f.WriteString(getTemInfo(v))
			f.Close()
		}
	}
	for _, m := range statics {
		if CheckFileIsExist(m) == false {
			dir, _ := path.Split(m)
			os.MkdirAll(dir, os.ModePerm)
			f, err := os.Create(m)
			if err != nil {
				fmt.Println(err.Error())
				logging.Fatal(err.Error())
			}
			f.WriteString(getStatics(m))
			f.Close()
		}
	}
}

var (
	FileTemplate = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>图库管理power by www.nbnat.com</title><meta name="renderer" content="webkit"><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"><link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css"><style>.layui-upload-img{width:92px;height:92px;margin:0 10px 10px 0}.layui-elip{overflow:hidden;white-space:nowrap;text-overflow:ellipsis;width:100px}</style></head><body style="padding:10px"><button type="button" class="layui-hide" id="test1"></button><div class="layui-fluid"><div id="fileManager" lay-filter="test"></div></div></body><script src="https://www.layuicdn.com/layui/layui.js"></script><script>layui.extend({ fileManager: "/static/layui_exts/fileManager/fileManager" });
    layui.use(["fileManager", "layer", "upload"], function () {
      var fileManager = layui.fileManager,
        $ = layui.$,
        upload = layui.upload,
        layer = layui.layer;
      $("title").html($("title").html() + " version:" + fileManager.v);
      var upIns = upload.render({
        elem: "#test1", //绑定元素
        url: "data.php?action=upload", //上传接口
        field: "file[]",
      });
      fileManager.render({
        elem: "#fileManager",
        method: "post",
        id: "fmTest",
        btn_upload: true,
        btn_create: true,
        icon_url: "ico/",
        url: "/file_json",
        thumb: {
          nopic: "https://www.layuicdn.com/layui/images/face/62.gif",
          width: 100,
          height: 100,
        },
        parseData: function (res) {
          let _res = [];
          _res.code = 0;
          _res.data = res.images;
          _res.count = res.count;
          return _res;
        },
        done: function (res, curr, count) {},
        page: {
          limit: 30,
          layout: ["count", "prev", "page", "next", "limit", "refresh", "skip"],
        },
        where: {},
      });
      fileManager.on("pic(test)", function (obj) {
        var data = obj.data;
        layer.alert(JSON.stringify(data), {
          title: "当前数据：",
        });
      });
      fileManager.on("uploadfile(test)", function (obj) {
        upIns.config.data = { path: obj.path };
        upIns.config.done = function (res) {
          fileManager.reload("fmTest");
        };
        var e = document.createEvent("MouseEvents");
        e.initEvent("click", true, true);
        document.getElementById("test1").dispatchEvent(e);
      });
      fileManager.on("new_dir(test)", function (obj) {
        layer.msg("没有这个功能哦");
      });
      $(document).on("click", ".picview", function (res) {
        let name = $(this).data("name");
        layui.data("_fm", { key: "_picview_name", value: name });
        layer.open({
          type: 2,
          area: ["620px", "600px"],
          content: ["pop.html", "no"],
        });
      });
      window.callback = function (res) {
        let name = layui.data("_fm")._picview_name;
        console.log(name);
        $(".picview[data-name=" + name + "]")
          .find("img")
          .attr("src", res.thumb);
        $(".picview[data-name=" + name + "]")
          .find("p")
          .text(res.name);
        layui.data("_fm", null);
      };
    });</script></html>`
)
