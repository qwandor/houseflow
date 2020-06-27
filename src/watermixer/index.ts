import express from 'express';
import fetch from 'node-fetch';
import { WaterRequestType } from '../types';

import { isAuthenticated } from '../auth';
import { sendMessage } from '../firebase';
import { getIp } from '../helpers';
import { WatermixerData } from '@gbaranski/types';

const url: string = 'http://192.168.1.120';
let isProcessing: boolean = false;
let watermixerData: WatermixerData;
export async function waterMixerHandleRequest(
  req: express.Request,
  res: express.Response,
  requestType: WaterRequestType,
) {
  if (!isAuthenticated(req.header('username') || '', req.header('password') || '')) {
    console.log(`${getIp(req)} with ${req.get('user-agent')} on ${requestType} not authenticated`);
    res.status(401).end();
    return;
  }
  console.log(`${getIp(req)} with ${req.get('user-agent')} on ${requestType} authenticated`);
  switch (requestType) {
    case WaterRequestType.GET_DATA:
      res.json(JSON.stringify(watermixerData));
      break;
    case WaterRequestType.START_MIXING:
      await res.status(await waterMixerFetchUrl(WaterRequestType.START_MIXING)).end();
      if (req.header('username') !== 'gbaranski') {
        sendMessage(req.header('username') || '', `watermixer${requestType}`);
      }
      break;
    default:
      res.status(500).end();
      break;
  }
}

export async function waterMixerFetchUrl(path: string): Promise<number> {
  isProcessing = true;
  let statusCode = 0;
  await fetch(url + path, {
    method: 'POST',
  })
    .then(data => {
      console.log('Success:', data.status);
      statusCode = data.status;
    })
    .catch(error => {
      console.error('Error:', error);
      statusCode = 503;
    });
  isProcessing = false;
  return statusCode;
}

export async function waterMixerFetchEspDataInterval() {
  if (isProcessing) {
    console.log('Connection overloaded');
    return;
  }
  isProcessing = true;
  fetch(url + WaterRequestType.GET_DATA)
    .then(res => res.json())
    .then(data => {
      watermixerData = data;
      console.log('Fetched watermixer data');
    })
    .finally(() => {
      isProcessing = false;
    });
}
