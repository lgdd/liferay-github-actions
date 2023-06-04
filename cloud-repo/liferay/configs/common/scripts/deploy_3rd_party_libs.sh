#!/usr/bin/env bash

TMPDIR=/tmp/deploy
DEPLOYDIR=/opt/liferay/deploy
mkdir -p $TMPDIR
### You can add more modules here ###
modules=(
    "https://github.com/peterrichards-lr/liferay-twilio-integration/releases/download/1.0.3/com.liferay.twilio.whatsapp.api-1.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-twilio-integration/releases/download/1.0.3/com.liferay.twilio.whatsapp.service-1.0.3.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.account.entry.creator-3.0.2.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.context.inspector-1.0.0.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.custom.field.updater-3.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.dynamic.data.mapping.form.action.outcome.evaluator-6.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.dynamic.data.mapping.form.extractor-6.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.dynamic.data.mapping.form.mailer-3.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.dynamic.data.mapping.form.object.extractor-2.0.0.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.dynamic.data.mapping.form.options.translator-4.0.0.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.dynamic.data.mapping.upload.processor-5.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.extensions.common-6.0.2.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.organisation.creator-3.0.3.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.user.account.creator-3.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.user.group.creator-1.0.0.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.user.group.roles.updater-1.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.user.group.updater-1.0.1.jar"
    "https://github.com/peterrichards-lr/liferay-workflow-extensions/releases/download/7.0.0/com.liferay.workflow.whatsapp-2.0.2.jar"
)
for module in ${modules[@]}; do
  (cd $TMPDIR && curl -L -O $module)
done
mv $TMPDIR/* $DEPLOYDIR