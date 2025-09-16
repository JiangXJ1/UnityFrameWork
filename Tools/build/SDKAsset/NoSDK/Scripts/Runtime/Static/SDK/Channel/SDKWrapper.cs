public class SDKWrapper : SDKWrapperBase
{
    public override SDKType SDKType => SDKType.None;

    public override void DownloadConfig()
    {
        //按需下载
    }
    public override bool IsConfigDownload()
    {
        return true;
    }

    public override string[] GetServerList()
    {
        return null;
    }

    public override string GetUserId()
    {
        return string.Empty;
    }

    public override bool GetGameMaintenance(out string title, out string message, out string bundleUrl, out long timeLeft, out string[] whiteUids)
    {
        title = string.Empty;
        message = string.Empty;
        bundleUrl = string.Empty;
        whiteUids = null;
        timeLeft = 0;
        return false;
    }

    public override bool GetServerCloseedInfo(out string title, out string message, out string[] whiteUids)
    {
        title = string.Empty;
        message = string.Empty;
        whiteUids = null;
        return false;
    }

    public override string GetSvnRevision()
    {
        return string.Empty;
    }

    public override string GetStoreUrl()
    {
        return string.Empty;
    }

    public override string GetSocialUrl()
    {
        return string.Empty;
    }
}

