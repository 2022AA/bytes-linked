// 返回结果的数据类型
export interface resultType {
    code: number;
    data: object;
    msg: string;
}

export interface itemType {
    file_addr: string;
    file_sha1: string;
    img_addr: string;
    file_name: string;
    like_cnt: number;
    username?: string;
    ar_tag?: boolean;
    file_id?: number;
    status?: number;
    avatar_url?: string;
    owner_uid?: number;
}
