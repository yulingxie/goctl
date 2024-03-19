package jenkins

const ServiceFoldTpl = `
<com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.708.ve61636eb_65a_5">
<actions/>
<description/>
<displayName>{{.name}}-service</displayName>
<properties/>
<folderViews class="com.cloudbees.hudson.plugins.folder.views.DefaultFolderViewHolder">
<views>
<hudson.model.AllView>
<owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
<name>All</name>
<filterExecutors>false</filterExecutors>
<filterQueue>false</filterQueue>
<properties class="hudson.model.View$PropertyList"/>
</hudson.model.AllView>
</views>
<tabBar class="hudson.views.DefaultViewsTabBar"/>
</folderViews>
<healthMetrics/>
<icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>
`

const GatewayFoldTpl = `
<com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.708.ve61636eb_65a_5">
<actions/>
<description/>
<displayName>{{.name}}</displayName>
<properties/>
<folderViews class="com.cloudbees.hudson.plugins.folder.views.DefaultFolderViewHolder">
<views>
<hudson.model.AllView>
<owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
<name>All</name>
<filterExecutors>false</filterExecutors>
<filterQueue>false</filterQueue>
<properties class="hudson.model.View$PropertyList"/>
</hudson.model.AllView>
</views>
<tabBar class="hudson.views.DefaultViewsTabBar"/>
</folderViews>
<healthMetrics/>
<icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>
`
