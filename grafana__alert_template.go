{{- /* Value  название для данных  чтобы значение передалось в сообщение внутри самого алерта*/ -}}
{{- /* [CustomURLName](CustomURL) создать в Labels чтобы создать рефералку в сообщении*/ -}}
{{- /* LablesExceptionRegex создать в Labels чтобы создать исключения в лейблах по значения ключа в сообщении*/ -}}
{{- /* CustomUrlPanel создать в Labels чтобы вставить ссылку взамен ссылки на панель*/ -}}
{{- /* CustomUrlDashboard создать в Labels чтобы вставить ссылку взамен ссылки на дашборд*/ -}}


{{ define "telegram_EAv2MarkDown" }}

  {{- range .Alerts -}}

    {{- if not (match "^Datasource(Error|NoData)$" .Labels.alertname) -}}

      {{- if (match "firing" .Status) -}}
        {{- print "🔥 *Firing*" -}}
      {{- end -}}

      {{- if and (match "resolved" .Status) (not (match "(^Updated$|^MissingSeries$)" .Annotations.grafana_state_reason)) -}}
        {{- print "🍀 *Resolved*" -}}
      {{- end -}}

      {{- if not (match "^$" .Annotations.grafana_state_reason)  -}}
        {{- if not (match "(^Updated$|^MissingSeries$)" .Annotations.grafana_state_reason) -}}
          {{-  template "boldnolnHead" join " " (stringSlice "| [" .Annotations.grafana_state_reason "]" ) -}}
        {{- else -}}
          {{- if  (match "^Updated$" .Annotations.grafana_state_reason) -}}
          {{- template "boldnolnHead" join " " (stringSlice "🆙 " .Annotations.grafana_state_reason ) -}}
          {{- end -}}
          {{- if  (match "^MissingSeries$" .Annotations.grafana_state_reason) -}}
          {{- template "boldnolnHead" join " " (stringSlice "♻ " .Annotations.grafana_state_reason ) -}}
          {{- end -}}
        {{- end -}}
      {{- end -}}

      {{- template "boldnolnHead" join "" (stringSlice " | " .Labels.alertname ) -}}
      
      {{- if .Labels.rulename -}}
        {{-  template "boldnolnHead" join "" (stringSlice " | " .Labels.rulename ) -}}
      {{- end -}}
    
    {{- else -}}

      {{- if (match "firing" .Status) -}}
        {{- print "❌ " -}}
      {{- end -}}

      {{- if (match "resolved" .Status)  -}}
        {{- print "✅ *Resolved* | " -}}
      {{- end -}}
      {{- template "boldnolnHead" .Labels.alertname  -}}
      
      {{- if .Labels.rulename -}}
        {{-  template "boldnolnHead" join "" (stringSlice " | " .Labels.rulename ) -}}
      {{- end -}}

    {{- end -}}

    {{- println "" -}}

    {{- println "\n*Time start:*" (.StartsAt | tz "Europe/Moscow" | date "2006.01.02 15:04:05") -}}
    
    {{- if eq .Status "resolved" -}}
      {{- println "*Time resolved:*" (.EndsAt | tz "Europe/Moscow" | date "2006.01.02 15:04:05") -}}
    {{- end -}}
    
    {{- $Unit := .Annotations.Unit -}}
    {{- $CustomURL := .Annotations.CustomURL -}}
    {{- $CustomURLName := .Annotations.CustomURLName -}}
    {{- $LablesExceptionRegex := .Annotations.LablesExceptionRegex -}}
    {{- $CustomUrlPanel := .Annotations.CustomUrlPanel -}}
    {{- $CustomUrlDashboard := .Annotations.CustomUrlDashboard -}}

    {{- template "annotations" (.Annotations.Remove (stringSlice "CustomUrlDashboard" "CustomUrlPanel" "Unit" "CustomURL" "CustomURLName" "LablesExceptionRegex" )) -}}

    {{- /* LABLES */ -}}

    {{- $exceptionLabels := "" -}}
    {{- if not (match "^DatasourceError$" .Labels.alertname) -}}{{- $exceptionLabels = "datasource" -}}{{- end -}}
    {{- $exceptionLabels = ( stringSlice "ref_id" "datasource_uid" "rulename" "silence_temp" "class" "rule_uid" "thread" "alertname" "item_key" "metrics" "metric" "type" "grafana_folder" "service" $exceptionLabels ) -}}

    {{- if len (.Labels.Remove $exceptionLabels) -}}
      {{- template "bold" "Labels:" -}}
      {{- range $key,$val := (.Labels.Remove $exceptionLabels) -}}

        {{- $key = reReplaceAll "([_*])" `$1$1\$1` $key -}}
        {{- $val = reReplaceAll "([_*])" `$1$1\$1` $val -}}

        {{- if $LablesExceptionRegex -}}
          {{- if not (match $LablesExceptionRegex $key) -}}
          
            {{- template "italic_key" $key -}}
            {{- println $val -}}

          {{- end -}}
        {{- else -}}

          {{- template "italic_key" $key -}}
          {{- println $val -}}

        {{- end -}}
      {{- end -}}
    {{- end -}}

    {{- /* END LABLES */ -}}   

    {{- if and $CustomURL $CustomURLName -}}
        {{- template "mark_ref" (stringSlice $CustomURLName $CustomURL) -}}
    {{- end -}}

    
    {{- if ( match "var='Value'" .ValueString ) -}}
      {{- println "" -}}
      {{- $value := print .Values.Value -}}
      {{- template "unit" (stringSlice .Status $value $Unit ) -}}
    {{- end -}}

    {{- if or .PanelURL .DashboardURL -}}
      {{- println "" -}}
    {{- end -}}

    {{- if .PanelURL -}}
      {{- if not $CustomUrlPanel -}}
        {{- template "mark_ref" (stringSlice "Panel" .PanelURL) -}}
      {{- else -}}
        {{- template "mark_ref" (stringSlice "Panel" $CustomUrlPanel) -}}
      {{- end -}}
      {{- if or .DashboardURL $CustomUrlDashboard -}}{{- print " | " -}}{{- end -}}
    {{- end -}}

    {{- if or .DashboardURL $CustomUrlDashboard -}}
      {{- if not $CustomUrlDashboard -}}
        {{- template "mark_ref" (stringSlice "Dashboard" .DashboardURL) -}}
      {{- else -}}
        {{- template "mark_ref" (stringSlice "Dashboard" $CustomUrlDashboard) -}}
      {{- end -}}

    {{- end -}}
    {{- println "" -}}

  {{- end -}}
{{- end -}}


{{- /* конец основной части  */ -}}

{{- define "bold" -}}
  {{- $text := reReplaceAll "(_)(_\\_)" `$1` . -}}
  {{- $text = reReplaceAll "([\\*]{2}\\\\)(\\*)" `$2$1$2$2` $text -}}
  {{- println (join "" (stringSlice "*" $text "*")) -}}
{{- end -}}

{{- define "boldnoln" -}}
  {{- $text := reReplaceAll "(_)(_\\_)" `$1` . -}}
  {{- $text = reReplaceAll "([\\*]{2}\\\\)(\\*)" `$2$1$2$2` $text -}}
  {{- print (join "" (stringSlice "*" $text "*")) -}}
{{- end -}}

{{- define "boldnolnHead" -}}
  {{- $text := reReplaceAll "([\\*])" `$1$1$1\$1$1` . -}}
  {{- print (join "" (stringSlice "*" $text "*")) -}}
{{- end -}}

{{- define "annotations" -}}
  
  {{- if len . -}}
    {{- template "bold" "Annotations:" -}}
    {{- range $key,$val := . -}}

        {{- $key = reReplaceAll "([_*])" `$1$1\$1` $key -}}
        {{- $val = reReplaceAll "([_*])" `$1$1\$1` $val -}}     

        {{- template "italic_key" $key -}}
        {{- println $val -}}
        
    {{- end -}}
  {{- end -}}

{{- end -}}

{{- define "italic_key" -}}

    {{- $var := reReplaceAll "([\\*]{2}\\\\)(\\*)" `$2` . -}}
    {{- $var = reReplaceAll "([_]{2}\\\\)(_)" `$2$1$2$2` $var -}}

    {{- print (join "" (stringSlice "\t- _" $var ":_ " )) -}}

{{- end -}}


{{- define "mark_ref" -}}

  {{- if eq (len .) 2 -}}
    {{- print "[" (index . 0) "](" (index . 1) ")" -}}
  {{- end -}}

{{- end -}}



{{- define "unit" -}}

    {{- if gt (len .) 1 -}}
      
      {{- $Status := (index . 0) -}}
      {{- $Value := (reReplaceAll "(\\d*)((\\.|,|\\D)\\d\\d)(\\d*)" "$1$2" (index . 1)) -}}
      {{- $Unit := (index . 2) -}}
      
      {{- if eq $Status "resolved" -}}
        {{- print "```Value:" " " "✅" " " $Value -}}
      {{- else -}}
        {{- print "```Value:" " " "⭕" " " $Value -}}
      {{- end -}}

      {{- if $Unit -}}
          {{- print " " $Unit " ```" -}}
      {{- else -}}
          {{- print " ```" -}}
      {{- end -}}

    {{- end -}}

{{- end -}}
