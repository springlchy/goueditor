<!DOCTYPE html>
<html>
    <head>
        <title>Beego</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    </head>
    <body>
        <form class="form-horizontal" method="POST" enctype="multipart/form-data">
          <div class="form-group">
        	<label for="stocks" class="col-sm-2 control-label">文章内容:</label>
         	<div class="col-sm-8">
          <script type="text/javascript" charset="utf-8" src="/plugins/ueditor/ueditor.config.js"></script>
        <script type="text/javascript" charset="utf-8" src="/plugins/ueditor/ueditor.all.min.js"> </script>

        <div style="text-align:left;height:500px;overflow:auto">
          <script id="editor" name="editor" type="text/plain" style="height:500px;"></script>
        </div>

        <script type="text/javascript">

          var ue = UE.getEditor('editor',{
            toolbars: [
              [
                'source','preview','drafts','undo','redo',
                'fontfamily','fontsize','bold','italic','underline', 'strikethrough', 'subscript',
                'fontborder', 'superscript', 'blockquote', 'horizontal', 'removeformat',
                'link', 'unlink', 'forecolor', 'backcolor'
              ],
              [
                'indent', 'justifyleft', 'justifyright', 'justifycenter', 'justifyjustify', 'rowspacingtop', 'rowspacingbottom',
                'lineheight', 'paragraph', 'pagebreak', 'spechars',

                'insertorderedlist', //有序列表
                'insertunorderedlist', //无序列表

                'map', //Baidu地图
                'gmap', //Google地图
                'insertvideo', //视频
                'insertframe', //插入Iframe

                'insertcode' //代码语言
              ],
              [
                'simpleupload', //单图上传
                'imagenone', //默认
                'imageleft', //左浮动
                'imageright', //右浮动

                'imagecenter', //居中

                'customstyle', //自定义标题
                'autotypeset', //自动排版

                'background' //背景

              ],
              [
                'inserttable', //插入表格
                'insertrow', //前插入行
                'insertcol', //前插入列
                'mergeright', //右合并单元格
                'mergedown', //下合并单元格
                'deleterow', //删除行
                'deletecol', //删除列
                'splittorows', //拆分成行
                'splittocols', //拆分成列
                'splittocells', //完全拆分单元格
                'deletecaption', //删除表格标题
                'inserttitle', //插入标题
                'mergecells', //合并多个单元格
                'deletetable', //删除表格
                'insertparagraphbeforetable', //"表格前插入行"
                'edittable', //表格属性
                'edittd' //单元格属性
              ]
            ],
            autoHeightEnabled: true,
            autoFloatEnabled: true
          });

        </script>


        <script>ue.ready(function(){ue.setContent('');});</script>
        <textarea type='text' name='contentTxt' id='editorTxt' hidden='hidden'></textarea>
        	</div>
          </div>

        	<div class="col-sm-offset-2 col-sm-8">
        		<button type="submit" class="btn btn-default">确 认</button>
        	</div>
    </form>
</body>
</html>