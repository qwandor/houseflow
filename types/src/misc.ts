import { Device, AnyDeviceData } from './device';
import { Client } from './client';

export interface DateTime {
  hour: number;
  minute: number;
  second: number;
}

export type State = {
  state: boolean;
};

export interface TempHistory {
  unixTime: number;
  temperature: number;
}

export interface RequestHistory {
  user: string;
  requestPath: string;
  unixTime: number;
  ip: string;
  userAgent: string;
  country: string;
}

export interface CurrentConnections {
  devices: Device.ActiveDevice[];
  clients: Client.ActiveUser[];
}

enum CloudTopics {
  DEVICE_DATA = 'device_data',
  DEVICE_DISCONNECT = 'device_disconnect',
  DEVICE_REQUEST = 'device_request',
}
