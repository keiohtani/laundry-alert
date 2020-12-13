#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { LaundryNotificationStack } from '../lib/laundry-notification-stack';

const app = new cdk.App();
new LaundryNotificationStack(app, 'LaundryNotificationStack');
