// package: teleport.terminal.v1
// file: v1/cluster.proto

import * as jspb from "google-protobuf";

export class Cluster extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getConnected(): boolean;
  setConnected(value: boolean): void;

  hasAcl(): boolean;
  clearAcl(): void;
  getAcl(): ACL | undefined;
  setAcl(value?: ACL): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Cluster.AsObject;
  static toObject(includeInstance: boolean, msg: Cluster): Cluster.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Cluster, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Cluster;
  static deserializeBinaryFromReader(message: Cluster, reader: jspb.BinaryReader): Cluster;
}

export namespace Cluster {
  export type AsObject = {
    name: string,
    connected: boolean,
    acl?: ACL.AsObject,
  }
}

export class Access extends jspb.Message {
  getList(): boolean;
  setList(value: boolean): void;

  getRead(): boolean;
  setRead(value: boolean): void;

  getEdit(): boolean;
  setEdit(value: boolean): void;

  getCreate(): boolean;
  setCreate(value: boolean): void;

  getDelete(): boolean;
  setDelete(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Access.AsObject;
  static toObject(includeInstance: boolean, msg: Access): Access.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Access, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Access;
  static deserializeBinaryFromReader(message: Access, reader: jspb.BinaryReader): Access;
}

export namespace Access {
  export type AsObject = {
    list: boolean,
    read: boolean,
    edit: boolean,
    create: boolean,
    pb_delete: boolean,
  }
}

export class ACL extends jspb.Message {
  hasSessions(): boolean;
  clearSessions(): void;
  getSessions(): Access | undefined;
  setSessions(value?: Access): void;

  hasAuthconnectors(): boolean;
  clearAuthconnectors(): void;
  getAuthconnectors(): Access | undefined;
  setAuthconnectors(value?: Access): void;

  hasRoles(): boolean;
  clearRoles(): void;
  getRoles(): Access | undefined;
  setRoles(value?: Access): void;

  hasUsers(): boolean;
  clearUsers(): void;
  getUsers(): Access | undefined;
  setUsers(value?: Access): void;

  hasTrustedclusters(): boolean;
  clearTrustedclusters(): void;
  getTrustedclusters(): Access | undefined;
  setTrustedclusters(value?: Access): void;

  hasEvents(): boolean;
  clearEvents(): void;
  getEvents(): Access | undefined;
  setEvents(value?: Access): void;

  hasTokens(): boolean;
  clearTokens(): void;
  getTokens(): Access | undefined;
  setTokens(value?: Access): void;

  hasServers(): boolean;
  clearServers(): void;
  getServers(): Access | undefined;
  setServers(value?: Access): void;

  hasAppservers(): boolean;
  clearAppservers(): void;
  getAppservers(): Access | undefined;
  setAppservers(value?: Access): void;

  hasDbservers(): boolean;
  clearDbservers(): void;
  getDbservers(): Access | undefined;
  setDbservers(value?: Access): void;

  hasKubeservers(): boolean;
  clearKubeservers(): void;
  getKubeservers(): Access | undefined;
  setKubeservers(value?: Access): void;

  clearSshloginsList(): void;
  getSshloginsList(): Array<string>;
  setSshloginsList(value: Array<string>): void;
  addSshlogins(value: string, index?: number): string;

  hasAccessrequests(): boolean;
  clearAccessrequests(): void;
  getAccessrequests(): Access | undefined;
  setAccessrequests(value?: Access): void;

  hasBilling(): boolean;
  clearBilling(): void;
  getBilling(): Access | undefined;
  setBilling(value?: Access): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ACL.AsObject;
  static toObject(includeInstance: boolean, msg: ACL): ACL.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ACL, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ACL;
  static deserializeBinaryFromReader(message: ACL, reader: jspb.BinaryReader): ACL;
}

export namespace ACL {
  export type AsObject = {
    sessions?: Access.AsObject,
    authconnectors?: Access.AsObject,
    roles?: Access.AsObject,
    users?: Access.AsObject,
    trustedclusters?: Access.AsObject,
    events?: Access.AsObject,
    tokens?: Access.AsObject,
    servers?: Access.AsObject,
    appservers?: Access.AsObject,
    dbservers?: Access.AsObject,
    kubeservers?: Access.AsObject,
    sshloginsList: Array<string>,
    accessrequests?: Access.AsObject,
    billing?: Access.AsObject,
  }
}

