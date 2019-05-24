CREATE TABLE IF NOT EXISTS `caiji_admin` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `username` varchar(20) NOT NULL COMMENT '管理员名称',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '管理员密码',
  `login_time` int(10) NOT NULL DEFAULT '0' COMMENT '登录时间',
  `login_num` int(11) NOT NULL DEFAULT '0' COMMENT '登录次数',
  `gid` tinyint(1) NOT NULL DEFAULT '0' COMMENT '权限类型',

  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='管理员表';

CREATE TABLE IF NOT EXISTS `caiji_link` (
  `link_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '索引id',
  `link_title` varchar(100) DEFAULT NULL COMMENT '标题',
  `link_url` varchar(100) DEFAULT NULL COMMENT '链接',
  `link_pic` varchar(100) DEFAULT NULL COMMENT '图片',
  `link_sort` tinyint(3) unsigned NOT NULL DEFAULT '255' COMMENT '排序',
  PRIMARY KEY (`link_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='友情链接表';


CREATE TABLE IF NOT EXISTS `caiji_posts` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `post_author` varchar(50) NOT NULL DEFAULT '0' COMMENT '作者',
  `post_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '发表时间',
  `post_content` longtext NOT NULL COMMENT '文章内容',
  `post_title` text NOT NULL COMMENT '文章标题',
  `post_source` text NOT NULL COMMENT '文章来源',
  `CssFile` text NOT NULL DEFAULT '' COMMENT 'css文件路径',
  `post_status` varchar(20) NOT NULL DEFAULT '已发表' COMMENT '文章状态 已发表/已删除/草稿',
  `comment_status` varchar(20) NOT NULL DEFAULT 'open' COMMENT '允许评论 open/close',
  `post_password` varchar(255) NOT NULL DEFAULT '' COMMENT '文章密码',
  `post_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `uuid` varchar(50) NOT NULL DEFAULT '',
  `comment_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '评论数',
  `look_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '浏览量',
  PRIMARY KEY (`ID`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='文章数据表';


CREATE TABLE IF NOT EXISTS `caiji_term_relationships` (
  `object_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `term_taxonomy_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`object_id`,`term_taxonomy_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='标签关系表';

CREATE TABLE IF NOT EXISTS `caiji_terms` (
  `term_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`term_id`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='标签表';

CREATE TABLE IF NOT EXISTS `caiji_comments` (
  `comment_ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `comment_postID` bigint(20) unsigned NOT NULL DEFAULT '1',
  `comment_uuid` varchar(50) NOT NULL DEFAULT '',
  `comment_author` tinytext NOT NULL,
  `comment_author_IP` int unsigned NOT NULL COMMENT '整数型IP',
  `comment_status` varchar(20) NOT NULL DEFAULT '已发表' COMMENT '评论状态 已发表/已删除/已屏蔽',
  `comment_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '发表时间',
  `comment_content` text NOT NULL,
  `comment_agent` varchar(255) NOT NULL DEFAULT '',
  `comment_parent` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '楼主评论id',
  PRIMARY KEY (`comment_ID`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='评论内容表';

CREATE TABLE IF NOT EXISTS `caiji_setting` (
  `ID` int(5) unsigned NOT NULL DEFAULT '0',
  `name` varchar(50) NOT NULL DEFAULT '',
  `value` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='设置表';
